package main

import (
	"fmt"
	"oengus-timers/sql"
	"oengus-timers/structs"
	"os"
	"os/signal"
	"time"
)

func StartTimers() {
	// Run everyting now!
	timerCallback()

	// tick every 5 minutes to check for updates
	ticker := time.NewTicker(5 * time.Minute)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	func() {
		for {
			select {
			case <-ticker.C:
				// Run everything on timer tick
				timerCallback()
			case <-sigint:
				ticker.Stop()
				return
			}
		}
	}()
}

func timerCallback() {
	go CheckSubmissionOpenClose()
}

func CheckSubmissionOpenClose() {
	marathonsToOpen, err1 := sql.FindMarathonsToOpenSubmissions()

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "(FindMarathonsToOpenSubmissions) QueryRow failed: %v\n", err1)
		return
	}

	openMarathonSubmissions(marathonsToOpen)

	marathonsToClose, err2 := sql.FindMarathonsToCloseSubmissions()

	if err2 != nil {
		fmt.Fprintf(os.Stderr, "(FindMarathonsToCloseSubmissions) QueryRow failed: %v\n", err2)
		return
	}

	closeMarathonSubmissions(marathonsToClose)
}

func extractIds(marathons []structs.Marathon) []string {
	v := make([]string, 0, len(marathons))

	for _, marathon := range marathons {
		v = append(v, marathon.Id)
	}

	return v
}

func openMarathonSubmissions(marathons []structs.Marathon) {
	// Open submissions
	// Enable edits for submissions

	marathonIds := extractIds(marathons)

	sql.OpenSubmission(marathonIds)
}

func closeMarathonSubmissions(marathons []structs.Marathon) {
	// Close submissions
	// Keep edits enabled

	marathonIds := extractIds(marathons)

	sql.CloseSubmission(marathonIds)
}
