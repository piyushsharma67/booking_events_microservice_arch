package models

import (
	"time"

	"github.com/piyushsharma67/events_booking/services/events_service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreateRequestToDocument(
	req *CreateEventRequest,
	organiserId string,
) (*EventDocument, error) {

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	rows := make([]SeatingRow, len(req.Rows))
	for i, r := range req.Rows {
		rows[i] = SeatingRow{
			RowLabel:  r.RowLabel,
			SeatCount: r.Seats,
		}
	}

	now := time.Now()

	return &EventDocument{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Location:    req.Location,
		StartTime:   startTime,
		EndTime:     endTime,
		Rows:        rows,
		Status:      utils.DRAFT,
		CreatedAt:   now,
		UpdatedAt:   now,
		OrganizerID: organiserId,
	}, nil
}
