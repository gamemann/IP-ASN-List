package webserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Info struct {
	token string
}

func (info *Info) handler(w http.ResponseWriter, r *http.Request) {
	// Check authorization token and header.
	if r.Header.Get("Authorization") != info.token {
		fmt.Fprintf(w, "Not authorized.")

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	filename := filepath.Base(r.URL.Path)

	// Try reading correct file.
	file, err := os.Open("public/" + filename)

	// Check for errors.
	if err != nil {
		fmt.Fprintf(w, "File not found.")

		w.WriteHeader(http.StatusNotFound)

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
		fmt.Fprintf(w, "Error reading file.")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	fmt.Fprintf(w, string(data))
}

func CreateWebServer(token string, port int) bool {
	// Create info struct and assign lists argument to lists member.
	info := &Info{token: token}

	// Create handler.
	http.HandleFunc("/", info.handler)

	// Create web server and have it listen on specific port.
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	// Check for errors.
	if err != nil {
		return false
	}

	return true
}
