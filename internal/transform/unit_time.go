package transform

import (
	"event-tracking-service/internal/handlers"
	"event-tracking-service/internal/models"
	"fmt"
	"sort"
)

func transformUnitTime(events []models.Event) string {
	if len(events) == 0 {
		return "no_events"
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})

	firstEvent := events[0]
	baseTime := firstEvent.CreatedAt

	// Format: "T0" for first event, "T1", "T2", etc. for subsequent events based on minute intervals
	timeUnits := make([]string, len(events))
	for i, event := range events {
		minutesDiff := int(event.CreatedAt.Sub(baseTime).Minutes())
		timeUnits[i] = fmt.Sprintf("T%d", minutesDiff)
	}

	// Return the unit time sequence
	return fmt.Sprintf("%s:%s", firstEvent.CaseID, timeUnits[len(timeUnits)-1])
}

func TransformUnitTime(caseID string) (string, error) {
	events, err := handlers.FetchEventsByCaseID(caseID)
	if err != nil {
		return "", err
	}

	unitTime := transformUnitTime(events)
	return unitTime, nil
}
