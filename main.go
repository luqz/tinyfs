package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	listen = flag.String("listen", ":8082", "listen address (e.g. :8082, 0.0.0.0:8082)")
	dir    = flag.String("dir", "uploads", "upload directory")
)

func main() {
	flag.Parse()

	uploadDir := *dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) { handleUploadPath(w, r, uploadDir) })
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) { handleUploadMultipart(w, r, uploadDir) })
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) { handleDownload(w, r, uploadDir) })
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handleList(w, r, uploadDir) })

	listenAddr := *listen
	log.Printf("Upload directory: %s", uploadDir)
	log.Printf("File server listening on http://%s", listenAddr)
	log.Printf("  Upload multipart:   curl -F \"file=@/path/to/file\" http://%s/upload", listenAddr)
	log.Printf("  Upload raw body:    curl --data-binary @/path/to/file http://%s/upload/filename", listenAddr)
	log.Printf("  Download:           curl -O http://%s/download/filename", listenAddr)
	log.Printf("  List files:         curl http://%s/", listenAddr)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func handleUploadMultipart(w http.ResponseWriter, r *http.Request, uploadDir string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	if filename == "" {
		http.Error(w, "filename is empty", http.StatusBadRequest)
		return
	}

	savePath := filepath.Join(uploadDir, filepath.Base(filename))
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save file: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK: uploaded %s (%d bytes)\n", filename, written)
	log.Printf("Uploaded: %s (%d bytes)", filename, written)
}

func handleUploadPath(w http.ResponseWriter, r *http.Request, uploadDir string) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/upload/")
	if filename == "" {
		http.Error(w, "filename required in path", http.StatusBadRequest)
		return
	}

	savePath := filepath.Join(uploadDir, filepath.Base(filename))
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save file: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK: uploaded %s (%d bytes)\n", filename, written)
	log.Printf("Uploaded: %s (%d bytes)", filename, written)
}

func handleDownload(w http.ResponseWriter, r *http.Request, uploadDir string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	filePath := filepath.Join(uploadDir, filepath.Base(filename))

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))
	http.ServeFile(w, r, filePath)
}

func handleList(w http.ResponseWriter, r *http.Request, uploadDir string) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "%-40s %10d bytes\n", entry.Name(), info.Size())
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "(empty)")
	}
}
