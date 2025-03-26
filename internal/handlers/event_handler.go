package handlers

import (
	"event-tracking-service/internal/database"
	"event-tracking-service/internal/models"
	"event-tracking-service/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	startTime := time.Now()
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		utils.LogError(fmt.Sprintf("Invalid event data: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate event_type
	if event.EventType != "start" && event.EventType != "end" {
		utils.LogError(fmt.Sprintf("Invalid event type: %s", event.EventType))
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_type must be either 'start' or 'end'"})
		return
	}

	// Parse the timestamp from the request
	if timestamp, err := time.Parse(time.RFC3339, event.Timestamp); err == nil {
		event.CreatedAt = timestamp
		event.UpdatedAt = timestamp
	} else {
		utils.LogError(fmt.Sprintf("Invalid timestamp format: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	// Get emoji for event type
	eventEmoji := utils.EventEmojis[event.EventName]
	if eventEmoji == "" {
		eventEmoji = "📎" // Default emoji for unknown event types
	}

	utils.LogInfo(fmt.Sprintf("Processing %s %s event for case %s", event.EventType, event.EventName, event.CaseID), eventEmoji)

	if err := database.DB.Create(&event).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to create event: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	duration := time.Since(startTime)
	utils.LogSuccess(fmt.Sprintf("Created %s %s event for case %s", event.EventType, event.EventName, event.CaseID))
	utils.LogAPI("POST", "/events", "201", duration)

	c.JSON(http.StatusCreated, event)
}

func GetAllEvents(c *gin.Context) {
	startTime := time.Now()
	var events []models.Event

	utils.LogInfo("Fetching all events", utils.StatusEmojis["database"])

	if err := database.DB.Find(&events).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve events: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	duration := time.Since(startTime)
	utils.LogSuccess(fmt.Sprintf("Retrieved %d events", len(events)))
	utils.LogAPI("GET", "/events", "200", duration)

	c.JSON(http.StatusOK, events)
}

func GetEventsByCaseID(c *gin.Context) {
	startTime := time.Now()
	caseID := c.Param("caseID")
	var events []models.Event

	utils.LogInfo(fmt.Sprintf("Fetching events for case %s", caseID), utils.StatusEmojis["database"])

	if err := database.DB.Where("case_id = ?", caseID).Find(&events).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve events for case %s: %v", caseID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	duration := time.Since(startTime)
	utils.LogSuccess(fmt.Sprintf("Retrieved %d events for case %s", len(events), caseID))
	utils.LogAPI("GET", fmt.Sprintf("/cases/%s/events", caseID), "200", duration)

	c.JSON(http.StatusOK, events)
}

func FetchEventsByCaseID(caseID string) ([]models.Event, error) {
	startTime := time.Now()
	var events []models.Event

	utils.LogInfo(fmt.Sprintf("Fetching events for case %s", caseID), utils.StatusEmojis["database"])

	if err := database.DB.Where("case_id = ?", caseID).Find(&events).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve events for case %s: %v", caseID, err))
		return nil, err
	}

	duration := time.Since(startTime)
	utils.LogSuccess(fmt.Sprintf("Retrieved %d events for case %s", len(events), caseID))
	utils.LogAPI("GET", fmt.Sprintf("/internal/cases/%s/events", caseID), "200", duration)

	return events, nil
}

func GetCaseMetrics(c *gin.Context) {
	startTime := time.Now()
	caseID := c.Param("caseID")
	var events []models.Event

	utils.LogInfo(fmt.Sprintf("Calculating metrics for case %s", caseID), utils.StatusEmojis["metrics"])

	if err := database.DB.Where("case_id = ?", caseID).Order("created_at").Find(&events).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve case metrics for %s: %v", caseID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve case metrics"})
		return
	}

	if len(events) == 0 {
		utils.LogWarning(fmt.Sprintf("No events found for case %s", caseID))
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this case"})
		return
	}

	firstEvent := events[0].CreatedAt
	lastEvent := events[len(events)-1].CreatedAt
	totalDuration := lastEvent.Sub(firstEvent)
	formattedDuration := utils.FormatDuration(totalDuration)

	utils.LogInfo(fmt.Sprintf("Case %s metrics: %d events over %s",
		caseID, len(events), formattedDuration),
		utils.StatusEmojis["time"])

	metrics := gin.H{
		"case_id":           caseID,
		"total_events":      len(events),
		"first_event_time":  firstEvent,
		"last_event_time":   lastEvent,
		"total_duration_ms": totalDuration.Milliseconds(),
		"duration_human":    formattedDuration,
	}

	duration := time.Since(startTime)
	utils.LogSuccess(fmt.Sprintf("Generated metrics for case %s", caseID))
	utils.LogAPI("GET", fmt.Sprintf("/cases/%s/metrics", caseID), "200", duration)

	c.JSON(http.StatusOK, metrics)
}
