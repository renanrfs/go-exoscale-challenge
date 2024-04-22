package main

import (
	"fmt"
	"time"
)

// UsageRecord represents a record of usage of resources.
type UsageRecord struct {
	Event     string
	Resource  string
	UUID      string
	Account   string
	Timestamp time.Time
}

// BillingStatement represents a billing statement for a resource.
type BillingStatement struct {
	Resource string
	UUID     string
	Account  string
	Duration int // Duration in minutes
}

// processUsage computes billing statements from a list of usage records.
func processUsage(events []UsageRecord) []BillingStatement {
	usageMap := make(map[string]UsageRecord)
	var statements []BillingStatement

	for _, event := range events {
		key := event.Account + "-" + event.UUID
		if event.Event == "create" {
			usageMap[key] = event
		} else if event.Event == "destroy" {
			if createEvent, exists := usageMap[key]; exists {
				duration := event.Timestamp.Sub(createEvent.Timestamp).Minutes()
				statement := BillingStatement{
					Resource: event.Resource,
					UUID:     event.UUID,
					Account:  event.Account,
					Duration: int(duration),
				}
				statements = append(statements, statement)
				delete(usageMap, key)
			}
		}
	}

	return statements
}

func main() {
	events := []UsageRecord{
		{
			Event:     "create",
			Resource:  "object",
			UUID:      "d8377d93-db71-488a-b894-54a962760bea",
			Account:   "ee12577c-983f-4729-a0e9-c5789a906c04",
			Timestamp: time.Date(2017, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Event:     "create",
			Resource:  "object",
			UUID:      "3e5d7420-0027-47d1-8ebd-5bd0e8078790",
			Account:   "9f3eb789-0c2c-4cf6-8861-d5afb6a207e5",
			Timestamp: time.Date(2017, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Event:     "destroy",
			Resource:  "object",
			UUID:      "3e5d7420-0027-47d1-8ebd-5bd0e8078790",
			Account:   "9f3eb789-0c2c-4cf6-8861-d5afb6a207e5",
			Timestamp: time.Date(2017, 3, 10, 2, 0, 0, 0, time.UTC),
		},
		{
			Event:     "destroy",
			Resource:  "object",
			UUID:      "d8377d93-db71-488a-b894-54a962760bea",
			Account:   "ee12577c-983f-4729-a0e9-c5789a906c04",
			Timestamp: time.Date(2017, 3, 10, 1, 0, 0, 0, time.UTC),
		},
	}

	statements := processUsage(events)
	for _, statement := range statements {
		fmt.Printf("Billing Statement: %+v\n", statement)
	}
}
