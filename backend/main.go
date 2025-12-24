package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wader/goutubedl"
	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// Basic MIME detection for audio formats YouTube provides
func mimeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case "m4a", "mp4":
		return "audio/mp4"
	case "webm":
		return "audio/webm"
	case "opus":
		return "audio/ogg" // Opus is often stored in Ogg
	case "mp3":
		return "audio/mpeg"
	default:
		return "application/octet-stream"
	}
}

// Strip unsafe characters from filenames
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "audio"
	}
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, name)
}

// Ensure the filename doesn't collide; if it exists, append a timestamp.
func uniquePath(dir, base, ext string) string {
	full := filepath.Join(dir, fmt.Sprintf("%s.%s", base, ext))
	if _, err := os.Stat(full); err != nil {
		// doesn't exist or other error — return as-is and let open handle perms
		return full
	}
	// add timestamp to avoid collision
	ts := time.Now().Format("20060102-150405")
	return filepath.Join(dir, fmt.Sprintf("%s-%s.%s", base, ts, ext))
}

func getMetaData(id string) (musicbrainzws2.Recording, error) {
	client := musicbrainzws2.NewClient(musicbrainzws2.AppInfo{
		Name:    "my-tool",
		Version: "1.0",
	})
	defer client.Close()
	ctx := context.Background()

	filter := musicbrainzws2.IncludesFilter{
		Includes: []string{
			"releases",       // release-list containing this recording
			"artist-credits", // artist credit info
			"isrcs",          // ISRCs
			"tags",           // tags
			"genres",         // genres
			"recording-rels", // recording relations (samples/remixes/etc.)
			"work-rels",      // linked works
			"artist-rels",    // artist relations
			"url-rels",       // external URLs
		},
	}

	return client.LookupRecording(ctx, mbtypes.MBID(id), filter)

}

type audioTags struct {
	Title  string
	Artist string
	Album  string
	Date   string
	Genre  string
}

func artistCreditString(credits musicbrainzws2.ArtistCredit) string {
	var b strings.Builder
	for _, c := range credits {
		name := c.Name
		if name == "" {
			name = c.Artist.Name
		}
		if name == "" {
			continue
		}
		b.WriteString(name)
		b.WriteString(c.JoinPhrase)
	}
	return strings.TrimSpace(b.String())
}

func uniqueNonEmpty(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func addID3TXXX(args *[]string, description, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	*args = append(*args, "-metadata", fmt.Sprintf("TXXX:%s=%s", description, value))
}

func buildAudioTags(rec musicbrainzws2.Recording) audioTags {
	var tags audioTags
	tags.Title = rec.Title
	tags.Artist = artistCreditString(rec.ArtistCredit)

	if len(rec.Releases) > 0 {
		tags.Album = rec.Releases[0].Title
		if tags.Date == "" {
			tags.Date = rec.Releases[0].Date.String()
		}
	}
	if tags.Date == "" {
		tags.Date = rec.FirstReleaseDate.String()
	}

	var genres []string
	for _, g := range rec.Genres {
		genres = append(genres, g.Name)
	}
	for _, g := range rec.UserGenres {
		genres = append(genres, g.Name)
	}
	for _, t := range rec.Tags {
		genres = append(genres, t.Name)
	}
	for _, t := range rec.UserTags {
		genres = append(genres, t.Name)
	}
	for _, release := range rec.Releases {
		for _, g := range release.Genres {
			genres = append(genres, g.Name)
		}
		for _, g := range release.UserGenres {
			genres = append(genres, g.Name)
		}
		for _, t := range release.Tags {
			genres = append(genres, t.Name)
		}
		for _, t := range release.UserTags {
			genres = append(genres, t.Name)
		}
	}
	genres = uniqueNonEmpty(genres)
	if len(genres) > 0 {
		tags.Genre = strings.Join(genres, ", ")
	}

	return tags
}

func transcodeAndTagToMP3(inputPath, outputPath string, rec musicbrainzws2.Recording) error {
	tags := buildAudioTags(rec)
	tmpPath := outputPath + ".tmp.mp3"
	args := []string{"-y", "-i", inputPath, "-map", "0:a:0", "-vn"}
	if tags.Title != "" {
		args = append(args, "-metadata", fmt.Sprintf("title=%s", tags.Title))
	}
	if tags.Artist != "" {
		args = append(args, "-metadata", fmt.Sprintf("artist=%s", tags.Artist))
	}
	if tags.Album != "" {
		args = append(args, "-metadata", fmt.Sprintf("album=%s", tags.Album))
	}
	if tags.Date != "" {
		args = append(args, "-metadata", fmt.Sprintf("date=%s", tags.Date))
	}
	if tags.Genre != "" {
		args = append(args, "-metadata", fmt.Sprintf("genre=%s", tags.Genre))
	}
	addID3TXXX(&args, "MusicBrainz Track Id", string(rec.ID))

	var artistIDs []string
	for _, credit := range rec.ArtistCredit {
		artistIDs = append(artistIDs, string(credit.Artist.ID))
	}
	for _, artistID := range uniqueNonEmpty(artistIDs) {
		addID3TXXX(&args, "MusicBrainz Artist Id", artistID)
	}

	if len(rec.Releases) > 0 {
		release := rec.Releases[0]
		addID3TXXX(&args, "MusicBrainz Release Id", string(release.ID))
		if release.ReleaseGroup != nil {
			addID3TXXX(&args, "MusicBrainz Release Group Id", string(release.ReleaseGroup.ID))
		}
	}
	args = append(args, "-c:a", "libmp3lame", "-q:a", "2", "-id3v2_version", "3", "-write_id3v1", "1", tmpPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("ffmpeg not found in PATH")
		}
		return err
	}

	return os.Rename(tmpPath, outputPath)
}

func downloadAudioHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	musicbrainzId := r.URL.Query().Get("musicbrainzid")

	if id == "" {
		log.Error("Missing Youtube ID")
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	var metaData musicbrainzws2.Recording
	if musicbrainzId != "" {
		tMetaData, err := getMetaData(musicbrainzId)
		if err != nil {
			log.Error("Error fetching Metadata from Musibrainz", "musicbrainzid", musicbrainzId)
		} else {
			log.Debug("Recording", "title", tMetaData.Title)
			log.Debug("Recording user genres", "genres", tMetaData.UserGenres)
			log.Debug("Recording user tags", "tags", tMetaData.UserTags)
			log.Debug("Recording artist credit", "artist_credit", tMetaData.ArtistCredit)
			log.Debug("Recording first release date", "date", tMetaData.FirstReleaseDate)
			log.Debug("Recording tags", "tags", tMetaData.Tags)
			// Print the titles of all found releases
			for i, release := range tMetaData.Releases {
				log.Debug("Release", "index", i, "id", release.ID, "title", release.Title)
				log.Debug("Release score", "score", release.Score)
				log.Debug("Release disambiguation", "disambiguation", release.Disambiguation)
				log.Debug("Release annotation", "annotation", release.Annotation)
				log.Debug("Release genres", "genres", release.Genres)
			}
		}

		if err == nil {
			metaData = tMetaData
		}
	}
	log.Debug("Metadata", "recording", metaData)

	log.Debug("Received request for id %s \n", id)
	ctx := r.Context()

	result, err := goutubedl.New(ctx, id, goutubedl.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rc, err := result.Download(ctx, "bestaudio") // Only audio stream

	if err != nil {
		log.Error("Error downloading from Youtube", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	// Determine filename + ext
	filename := "audio"
	if t := sanitizeFilename(result.Info.Title); t != "" {
		filename = t
	}
	ext := result.Info.Format.Ext
	if ext == "" {
		ext = "m4a"
	}
	mime := mimeFromExt("mp3")

	dir := "./downloads"
	if err := ensureDir(dir); err != nil {
		log.Error("Error creating output directory", "err", err)
		http.Error(w, "failed to create output directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	inputPath := uniquePath(dir, filename, ext)
	outputPath := uniquePath(dir, filename, "mp3")

	f, err := os.OpenFile(inputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Error("Error creating file", "err", err)
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(f, rc)
	if err != nil {
		_ = f.Close()
		log.Error("Error writing file", "err", err)
		http.Error(w, "failed while writing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := f.Close(); err != nil {
		log.Error("Error closing file", "err", err)
		http.Error(w, "failed while closing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := transcodeAndTagToMP3(inputPath, outputPath, metaData); err != nil {
		log.Error("Error transcoding/tagging audio file", "err", err)
		http.Error(w, "failed while tagging file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.Remove(inputPath); err != nil {
		log.Error("Error removing source file", "err", err)
	}

	fmt.Fprintf(w, "Success! Downloaded %s (%s)", filepath.Base(outputPath), mime)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/download", downloadAudioHandler).Methods(http.MethodGet, http.MethodOptions)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Info("Server running on http://localhost:3333")
	log.Error(http.ListenAndServe(":3333", cors(r)))
}
