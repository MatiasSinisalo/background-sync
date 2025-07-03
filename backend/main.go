package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var lastUpdate int64
var nextIndex int
var currentPath string

func updatePath() {

	dir := "./static"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	size := 0
	for _, entry := range entries {

		if !entry.IsDir() {
			//fmt.Printf("file: %s\n", entry.Name())
			size++
		}
	}

	searchIndex := nextIndex

	count := 0
	for _, entry := range entries {

		if !entry.IsDir() {
			//fmt.Printf("!!!changed background!!!\n")
			//fmt.Printf("file: %s\n", entry.Name())
			if count >= searchIndex {
				currentPath = filepath.Join(dir, entry.Name())
				nextIndex++
				if nextIndex >= size {
					nextIndex = 0
				}

				return
			}

			count++

		}
	}

	//fmt.Printf("Number of files in %s: %d\n", dir, count)
	return
}

func imageHandler(w http.ResponseWriter, r *http.Request) {

	seconds := time.Now().Unix()
	difference := seconds - lastUpdate
	fmt.Printf("%d, size %d\n", nextIndex, currentPath)
	if difference > 3 {
		updatePath()
		lastUpdate = seconds
	}

	//fmt.Printf("current image path: %s\n, and next image index: %d", currentPath, nextIndex)

	file, err := os.Open(currentPath)
	if err != nil {
		http.Error(w, "Image not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the appropriate content type
	w.Header().Set("Content-Type", "image/jpeg") // or image/png, etc.

	info, _ := file.Stat()
	// Serve the file content
	http.ServeContent(w, r, filepath.Base(currentPath), info.ModTime(), file)
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	//create and save the file
	dst, err := os.Create("./static/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully!", handler.Filename)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := `
        <!DOCTYPE html>
        <html>
        <head><title>Upload File</title></head>
        <body>
            <h1>Upload a File</h1>
            <form method="POST" enctype="multipart/form-data" action="/dashboard">
                <input type="file" name="uploadfile" />
                <input type="submit" value="Upload" />
            </form>
        </body>
        </html>`
		t := template.Must(template.New("upload").Parse(tmpl))
		t.Execute(w, nil)

	case "POST":
		imageUploadHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	lastUpdate = time.Now().Unix()
	currentPath = "./static/wallpaper2.jpeg"
	nextIndex = 0

	log.Println("server started..")
	http.HandleFunc("/api/wallpaper", imageHandler)
	http.HandleFunc("/api/wallpaper/add", imageUploadHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
