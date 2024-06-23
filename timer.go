package main

import (
	"fmt"
	"log"
	"oengus-timers/sql"
	"oengus-timers/structs"
	"oengus-timers/utils"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func StartTimers() {
	// Run everything now!
	timerCallback()

	timerInterval, err := strconv.Atoi(utils.GetEnv("TIMER_INTERVAL_MINUTES", "5"))

	if err != nil {
		panic(err)
	}

	// tick every 5 minutes to check for updates
	ticker := time.NewTicker(time.Duration(timerInterval) * time.Minute)

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
	if len(marathons) == 0 {
		log.Println("No submissions to open!")
		return
	}

	marathonIds := extractIds(marathons)

	log.Println("Opening submissions for", marathonIds)

	sql.OpenSubmission(marathonIds)
}

func closeMarathonSubmissions(marathons []structs.Marathon) {
	if len(marathons) == 0 {
		log.Println("No submissions to close!")
		return
	}

	marathonIds := extractIds(marathons)

	log.Println("Closing submissions for", marathonIds)

	sql.CloseSubmission(marathonIds)
}
