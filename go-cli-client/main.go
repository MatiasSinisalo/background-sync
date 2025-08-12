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

func runCommand(cmdTemplate string) error {
	cmdStr := fmt.Sprintf(cmdTemplate, savePath)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadAndUpdateWallpaper(commandTemplate string) {
	fmt.Println("Updating wallpaper...")
	if err := downloadImage(); err != nil {
		fmt.Println("Download error:", err)
		return
	}
	if err := runCommand(commandTemplate); err != nil {
		fmt.Println("Wallpaper set error:", err)
	} else {
		fmt.Println("Wallpaper updated!")
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: wallpaper-sync '<command>'")
		fmt.Println("Example: wallpaper-sync 'gsettings set org.gnome.desktop.background picture-uri file://%s'")
		os.Exit(1)
	}

	commandTemplate := os.Args[1]

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	fmt.Println("Wallpaper sync started...")

	for {
		select {
		case <-ticker.C:
			downloadAndUpdateWallpaper(commandTemplate)
			continue
		}
	}
}
