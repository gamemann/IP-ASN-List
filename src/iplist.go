package main

import (
	"./config"
	"./list"
	"./webserver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func UpdateEverything(lists *list.Lists, timer *time.Ticker) {
	destroy := make(chan struct{})

	for {
		select {
		case <-timer.C:
			// First retreive new lists
			*lists = list.GetLists()

			// Update the public files.
			list.UpdateLists(lists)
		case <-destroy:
			timer.Stop()

			return
		}
	}
}

func main() {
	// Initialize config.
	cfg := config.Config{}

	// Read config file.
	config.ReadConfig(&cfg, "settings.conf")

	// Create lists variable.
	var lists list.Lists

	// Get all initial lists.
	lists = list.GetLists()

	// Update the public files.
	list.UpdateLists(&lists)

	// Create timer that'll update the lists every x seconds.
	ticker := time.NewTicker(time.Duration(cfg.UpdateTime) * time.Second)
	go UpdateEverything(&lists, ticker)

	// Create web server.
	log.Fatal(webserver.CreateWebServer(cfg.Token, cfg.Port))

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT)

	x := 0

	// Create a loop so the program doesn't exit. Look for signals and if SIGINT, stop the program.
	for x < 1 {
		kill := false
		s := <-sigc

		switch s {
		case os.Interrupt:
			kill = true
		}

		if kill {
			break
		}

		// Sleep every second to avoid unnecessary CPU consumption.
		time.Sleep(time.Duration(1) * time.Second)
	}
}
