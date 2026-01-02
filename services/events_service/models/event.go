package models

import "time"

type FetchSource string

const (
	SourceCache FetchSource = "cache"
	SourceDB    FetchSource = "db"
)

type SeatingRowResponse struct {
	RowLabel string `json:"row"`
	Seats    int    `json:"seats"`
}


type EventBase struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	ImageURL    string              `json:"image_url"`
	Location    string              `json:"location"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Rows        []SeatingRowResponse `json:"rows"`
}

type CreateEventRequest struct {
	EventBase
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	ImageURL    string              `json:"image_url"`
	Location    string              `json:"location" binding:"required"`
	StartTime   string              `json:"start_time" binding:"required"` // RFC3339
	EndTime     string              `json:"end_time" binding:"required"`
	Rows        []SeatingRowRequest `json:"rows" binding:"required,min=1"`
}

type SeatingRowRequest struct {
	RowLabel string `json:"row" binding:"required"`
	Seats    int    `json:"seats" binding:"required,min=1"`
}

type GetEventDetailtRequest struct {
    EventID string `form:"event_id"`
}

type EventDetailResponse struct{
	EventBase
}

type EventDetailResult struct {
	Data   *EventDetailResponse
	Source FetchSource
}