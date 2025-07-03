package main

import (
	"log"
	"os"
	"path/filepath"
	"fmt"
	"net/http"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
    imagePath := "./static/wallpaper1.jpeg" // Replace with your image path
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
