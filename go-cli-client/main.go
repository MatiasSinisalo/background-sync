package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	url      = "http://localhost:8080/api/wallpaper"
	savePath = "/tmp/current_wallpaper.jpg"
)

func downloadImage() error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func setWallpaper(path string) error {
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+path)
	return cmd.Run()
}

func downloadAndUpdateWallpaper() {
	fmt.Println("Updating wallpaper...")
	if err := downloadImage(); err != nil {
		fmt.Println("Download error:", err)
		return
	}
	if err := setWallpaper(savePath); err != nil {
		fmt.Println("Wallpaper set error:", err)
	} else {
		fmt.Println("Wallpaper updated!")
	}
}

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Wallpaper sync started...")

	for {
		select {
		case <-ticker.C:
			downloadAndUpdateWallpaper()
			continue
		}
	}
}
