package webserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Info struct {
	token string
}

func (info *Info) handler(w http.ResponseWriter, r *http.Request) {
	// Check authorization token and header.
	if r.Header.Get("Authorization") != info.token {
		w.WriteHeader(http.StatusUnauthorized)

		fmt.Fprintf(w, "Not authorized.")

		return
	}

	filename := filepath.Base(r.URL.Path)

	// Try reading correct file.
	file, err := os.Open("public/" + filename)

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprintf(w, "File not found.")

		return
	}

	// Create stat.
	stat, _ := file.Stat()

	// Make data.
	data := make([]byte, stat.Size())

	// Read data.
	_, err = file.Read(data)

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "Error reading file.")

		return
	}

	fmt.Fprintf(w, string(data))
}

func CreateWebServer(token string, port int) bool {
	// Create info struct and assign lists argument to lists member.
	info := &Info{token: token}

	// Create HTTP server.
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      http.HandlerFunc(info.handler),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	// Create web server and have it listen on specific port.
	err := srv.ListenAndServe()

	// Check for errors.
	if err != nil {
		return false
	}

	return true
}
