package rabbitmq

import (
	"encoding/json"
	"oengus-timers/structs"
	"time"
)

type rmqJsonSubmissionStatus struct {
	Open         bool   `json:"open"`
	MarathonName string `json:"marathon_name"`
	ClosesAt     string `json:"closes_at"`
}

type rmqJson struct {
	Event            string                  `json:"event"`
	Url              string                  `json:"url"`
	SubmissionStatus rmqJsonSubmissionStatus `json:"submission_status"`
}

func SendSubmissionsOpenEvents(marathons []structs.Marathon) {
	for _, marathon := range marathons {
		jsonData := rmqJson{
			Event: "SUBMISSION_OPEN_STATUS_CHANGED",
			Url:   marathon.Webhook,
			SubmissionStatus: rmqJsonSubmissionStatus{
				Open:         true,
				MarathonName: marathon.Name,
				ClosesAt:     marathon.SubmissionsEndDate.Format(time.RFC3339),
			},
		}

		openedEventJson, _ := json.Marshal(jsonData)
		PublishBotMessage(string(openedEventJson))
	}
}

func SendSubmissionsClosedEvents(marathons []structs.Marathon) {
	for _, marathon := range marathons {
		jsonData := rmqJson{
			Event: "SUBMISSION_OPEN_STATUS_CHANGED",
			Url:   marathon.Webhook,
			SubmissionStatus: rmqJsonSubmissionStatus{
				Open:         false,
				MarathonName: marathon.Name,
				ClosesAt:     marathon.SubmissionsEndDate.Format(time.RFC3339),
			},
		}

		closedEventJson, _ := json.Marshal(jsonData)
		PublishBotMessage(string(closedEventJson))
	}
}
