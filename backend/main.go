package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wader/goutubedl"
)

// --- Config ---

// AUDIO_OUTPUT_DIR controls where files are written. Falls back to ./downloads if unset.
func outputDir() string {
	if v := os.Getenv("AUDIO_OUTPUT_DIR"); strings.TrimSpace(v) != "" {
		return v
	}
	return "./downloads"
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// --- Utilities ---

// Sets proper Content-Type and Content-Disposition based on goutubedl.Info (kept in case you want it for other endpoints)
func setDownloadHeaders(w http.ResponseWriter, info *goutubedl.Info) {
	filename := "audio"
	ext := "m4a"

	if info != nil {
		if info.Title != "" {
			filename = sanitizeFilename(info.Title)
		}
		if info.Format.Ext != "" {
			ext = info.Format.Ext
		}
	}

	mime := mimeFromExt(ext)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s.%s"`, filename, ext),
	)

	// ⚠️ Don't set Content-Length unless it's exact — YouTube sometimes omits it
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

// --- HTTP ---

type saveResponse struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	MIME     string `json:"mime"`
	Size     int64  `json:"size"`
}

func downloadAudioHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request for id %s \n", id)
	ctx := r.Context()

	result, err := goutubedl.New(ctx, id, goutubedl.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rc, err := result.Download(ctx, "bestaudio") // Only audio stream
	if err != nil {
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

	// Ensure output directory exists
	dir := outputDir()
	if err := ensureDir(dir); err != nil {
		http.Error(w, "failed to create output directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build a unique file path
	outPath := uniquePath(dir, filename, ext)

	// Create destination file
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy stream to file
	n, err := io.Copy(f, rc)
	if err != nil {
		log.Printf("write error: %v", err)
		http.Error(w, "failed while writing file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Done! Sending shit")
	// Respond with JSON about the saved file
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saveResponse{
		Path:     outPath,
		Filename: filepath.Base(outPath),
		MIME:     mime,
		Size:     n,
	})
}

func main() {
	log.Printf("AUDIO_OUTPUT_DIR=%q (default ./downloads if empty)", outputDir())

	r := mux.NewRouter()
	r.HandleFunc("/download", downloadAudioHandler).Methods(http.MethodGet, http.MethodOptions)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	fmt.Println("Server running on http://localhost:3333")
	log.Fatal(http.ListenAndServe(":3333", cors(r)))
}
