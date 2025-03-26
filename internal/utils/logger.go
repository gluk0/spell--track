package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Emoji constants for different event types
var EventEmojis = map[string]string{
	"document_review":    "📄",
	"client_meeting":     "👥",
	"case_update":        "📝",
	"file_upload":        "📤",
	"payment_processing": "💰",
	"status_change":      "🔄",
	"note_added":         "📌",
	"email_sent":         "📧",
	"phone_call":         "📞",
	"document_signing":   "✍️",
}

// Status emojis for different log types
var StatusEmojis = map[string]string{
	"success":    "✅",
	"error":      "❌",
	"info":       "ℹ️",
	"warning":    "⚠️",
	"api":        "🌐",
	"database":   "🗄️",
	"metrics":    "📊",
	"time":       "⏱️",
	"processing": "⚙️",
}

// Color functions
var (
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	timeColor    = color.New(color.FgMagenta).SprintFunc()
	apiColor     = color.New(color.FgBlue).SprintFunc()
)

// LogInfo prints an info message with timestamp and emoji
func LogInfo(message string, emoji string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("%s %s | %s\n", emoji, timeColor(timestamp), infoColor(message))
}

// LogSuccess prints a success message with timestamp and emoji
func LogSuccess(message string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("%s %s | %s\n", StatusEmojis["success"], timeColor(timestamp), successColor(message))
}

// LogError prints an error message with timestamp and emoji
func LogError(message string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("%s %s | %s\n", StatusEmojis["error"], timeColor(timestamp), errorColor(message))
}

// LogWarning prints a warning message with timestamp and emoji
func LogWarning(message string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("%s %s | %s\n", StatusEmojis["warning"], timeColor(timestamp), warnColor(message))
}

// LogAPI prints an API-related message with timestamp and emoji
func LogAPI(method, path, status string, duration time.Duration) {
	timestamp := time.Now().Format("15:04:05.000")
	durationStr := fmt.Sprintf("%.2fms", float64(duration.Microseconds())/1000.0)
	message := fmt.Sprintf("%s %s [%s] (%s)", method, path, status, durationStr)
	fmt.Printf("%s %s | %s\n", StatusEmojis["api"], timeColor(timestamp), apiColor(message))
}

// FormatDuration formats a duration with appropriate units
func FormatDuration(d time.Duration) string {
	if d.Hours() > 1 {
		return fmt.Sprintf("%.1f hours", d.Hours())
	} else if d.Minutes() > 1 {
		return fmt.Sprintf("%.1f minutes", d.Minutes())
	} else if d.Seconds() > 1 {
		return fmt.Sprintf("%.1f seconds", d.Seconds())
	}
	return fmt.Sprintf("%d milliseconds", d.Milliseconds())
}
