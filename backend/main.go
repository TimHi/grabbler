package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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
			log.Error("Error fetching Metadata from Musibrainz", "id", id)
		}
		if err == nil {
			fmt.Println("Recording:")
			fmt.Println(tMetaData.Title)
			fmt.Println(tMetaData.Genres)
			fmt.Println(tMetaData.Annotation)
			fmt.Println(tMetaData.ArtistCredit)
			fmt.Println(tMetaData.FirstReleaseDate)
			fmt.Println(tMetaData.Tags)
			// Print the titles of all found releases
			for i, release := range tMetaData.Releases {
				fmt.Printf("Release %d \n", i)
				fmt.Println(release.ID)
				fmt.Println(release.Title)
				fmt.Println(release.Score)
				fmt.Println(release.Disambiguation)
				fmt.Println(release.Annotation)
				fmt.Println(release.Genres)
			}
		} else {
			log.Fatal(err)
		}

		metaData = tMetaData
	}
	fmt.Println(metaData)

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
	mime := mimeFromExt(ext)

	dir := "./downloads"
	if err := ensureDir(dir); err != nil {
		log.Error("Error creating output directory", "err", err)
		http.Error(w, "failed to create output directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outPath := uniquePath(dir, filename, ext)

	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Error("Error creating file", "err", err)
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, rc)
	if err != nil {
		log.Error("Error writing file", "err", err)
		http.Error(w, "failed while writing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success! Downloaded %s.%s", filename, mime)
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
