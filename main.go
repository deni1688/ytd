package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rylio/ytdl"
)

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/convert", convertRoute).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":4343",
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	fmt.Println("Server starting on port 4343. Started at: " + time.Now().Format(time.RFC3339))

	log.Fatal(srv.ListenAndServe())
}

func convertRoute(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	mp3, err := createMP3(query.Get("url"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.ServeFile(w, r, mp3)
	os.Remove(mp3)
}

func createMP3(url string) (string, error) {
	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		fmt.Println("Failed to get video info", err.Error())
		return "", errors.New("Failed to get video info")
	}

	from, to := getFromTo()

	err = os.Remove(to)
	if err == nil {
		fmt.Println("Creating new file")
	}

	downloadFile(vid, from)

	err = ffmpegConvert(from, to)
	if err != nil {
		fmt.Println("ffmpeg failed to convert file", err.Error())
		return to, err
	}

	os.Remove(from)

	return to, nil
}

func ffmpegConvert(from, to string) error {
	fmt.Println("Download complete")

	cmd := exec.Command("ffmpeg", "-i", from, "-map", "0:a:0", "-b:a", "96k", to)
	if err := cmd.Run(); err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		return err
	}
	fmt.Printf("New file %s created. Thank you come again!\n", to)

	return nil
}

func downloadFile(vid *ytdl.VideoInfo, from string) {
	file, _ := os.Create(from)

	vid.Download(vid.Formats[0], file)
	file.Close()
}

func getFromTo() (string, string) {
	filename := uuid.New().String()

	from := "./downloads/" + filename + ".mp4"
	to := "./downloads/" + filename + ".mp3"

	return from, to
}
