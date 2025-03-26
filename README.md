# Event Tracking Microservice

This microservice tracks individual touchpoints between case managers and their interactions with Merlin. It provides insights into the time taken for each touchpoint and the total processing time for cases.

## Features

- Track individual events with case IDs
- Store event metadata in JSON format
- Calculate time metrics for case processing
- RESTful API endpoints
- Containerized deployment with Docker

## API Endpoints

### 1. Create Event
```http
POST /events
Content-Type: application/json

{
    "case_id": "CASE123",
    "event_name": "document_review",
    "event_type": "start",
    "metadata": {
        "document_type": "invoice",
        "status": "completed"
    }
}
```

### 2. Get Events by Case ID
```http
GET /cases/{caseID}/events
```

### 3. Get Case Metrics
```http
GET /cases/{caseID}/metrics
```

## Setup and Running

1. Clone the repository
2. Make sure you have Docker and Docker Compose installed
3. Run the service:
   ```bash
   docker-compose up --build
   ```

The service will be available at `http://localhost:8080`

## Environment Variables

- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: event_tracking)
- `DB_PORT`: Database port (default: 5432)

## Database Schema

### Events Table
- `id`: Primary key
- `case_id`: Case reference ID
- `event_name`: Name of the event
- `event_type`: Type of event ("start" or "end")
- `metadata`: JSON field for additional event data
- `created_at`: Timestamp of event creation
- `updated_at`: Timestamp of last update

### Cases Table
- `case_id`: Primary key
- `status`: Current case status
- `created_at`: Timestamp of case creation
- `updated_at`: Timestamp of last update
- `completed_at`: Timestamp of case completion 