package main

import (
	"log"

	"event-tracking-service/internal/database"
	"event-tracking-service/internal/handlers"
	"event-tracking-service/internal/middleware"
	"event-tracking-service/internal/utils"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

const banner = `          
                            88 88         88                              88       
.d8888b.  88d888b. .d8888b. 88 88       d8888P 88d888b. .d8888b. .d8888b. 88  .dP  
Y8ooooo.  88'  '88 88ooood8 88 88         88   88'  '88 88'  '88 88'  '"" 88888"   
      88  88.  .88 88.  ... 88 88         88   88       88.  .88 88.  ... 88  '8b. 
'88888P'  88Y888P' '88888P' dP dP         dP   dP       '88888P8 '88888P' dP   'YP 
ooooooooo~88~oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo

             ===> An onlyPlans creation <===                                                 

`

func printBanner() {
	// Create a bright cyan color for the banner
	bannerColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	println("\n" + bannerColor(banner))

	// Print a separator line
	separatorColor := color.New(color.FgMagenta).SprintFunc()
	println(separatorColor("═══════════════════════════════════════════════════════════════════════════════"))
}

func main() {
	// Print the banner
	printBanner()

	utils.LogInfo("Starting server...", "🚀")
	database.InitDB()
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	r.POST("/events", handlers.CreateEvent)
	r.GET("/events", handlers.GetAllEvents)
	r.GET("/cases/:caseID/events", handlers.GetEventsByCaseID)
	r.GET("/cases/:caseID/metrics", handlers.GetCaseMetrics)

	utils.LogInfo("Server is ready to accept connections", "✨")
	if err := r.Run(":8080"); err != nil {
		utils.LogError("Failed to start server: " + err.Error())
		log.Fatal(err)
	}
}
