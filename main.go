package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kkdai/youtube/v2"
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
	handleErr(os.Remove(mp3))
}

func createMP3(url string) (string, error) {
	client := youtube.Client{}
	vid, err := client.GetVideo(url)
	if err != nil {
		fmt.Println("Failed to get video info", err.Error())
		return "", errors.New("Failed to get video info")
	}

	from, to := getFromTo()

	handleErr(os.Remove(to))

	downloadFile(vid, from, &client)

	err = ffmpegConvert(from, to)
	if err != nil {
		fmt.Println("ffmpeg failed to convert file", err.Error())
		return to, err
	}

	return to, os.Remove(from)
}

func ffmpegConvert(from, to string) error {
	fmt.Println("Download complete")

	err := exec.Command("ffmpeg", "-i", from, "-map", "0:a:0", "-b:a", "96k", to).Run()
	handleErr(err)

	fmt.Printf("New file %s created. Thank you come again!\n", to)

	return nil
}

func downloadFile(vid *youtube.Video, from string, c *youtube.Client) {
	resp, err := c.GetStream(vid, &vid.Formats[0])
	handleErr(err)

	defer resp.Body.Close()

	file, err := os.Create(from)
	handleErr(err)

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	handleErr(err)
}

func getFromTo() (string, string) {
	filename := uuid.New().String()

	from := "./downloads/" + filename + ".mp4"
	to := "./downloads/" + filename + ".mp3"

	return from, to
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
