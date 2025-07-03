package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getPath() string {
	seconds := time.Now().Unix()
	if seconds%2 == 0 {
		return "./static/wallpaper1.jpeg"
	} else {
		return "./static/wallpaper2.jpeg"
	}

}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	imagePath := getPath()
	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Image not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the appropriate content type
	w.Header().Set("Content-Type", "image/jpeg") // or image/png, etc.

	info, _ := file.Stat()
	// Serve the file content
	http.ServeContent(w, r, filepath.Base(imagePath), info.ModTime(), file)
}

func main() {
	log.Println("server started..")
	http.HandleFunc("/api/wallpaper", imageHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
