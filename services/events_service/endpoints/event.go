package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/piyushsharma67/events_booking/services/events_service/models"
	"github.com/piyushsharma67/events_booking/services/events_service/que"
	"github.com/piyushsharma67/events_booking/services/events_service/service"
)

func GenerateEvent(srv *service.EventService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		event := request.(*models.CreateEventRequest)
		// Get organizer ID from context
		organizerID, ok := ctx.Value("user_id").(string)
		if !ok || organizerID == "" {
			return nil, errors.New("organizer ID not found in context")
		}

		genEvent, err := srv.CreateEvent(ctx, event, organizerID)

		if err != nil {
			return nil, err
		}

		// 2️⃣ Build message struct
		msg := que.GenerateSeatsMessage{
			EventID: genEvent.ID.Hex(),
		}

		for _, r := range event.Rows {
			msg.SeatLayout.Rows = append(msg.SeatLayout.Rows, struct {
				Row   string `json:"row"`
				Count int    `json:"count"`
			}{
				Row:   r.RowLabel,
				Count: r.Seats,
			})
		}

		// 3️⃣ Marshal
		payload, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		err = srv.Publisher.Publish(ctx, payload)
		if err != nil {
			return nil, err
		}

		return genEvent, nil
	}
}

func GetEventDetails(srv *service.EventService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		event := request.(*models.GetEventDetailtRequest)
		fmt.Println("EventID:", event.EventID)
		// Get organizer ID from context
		organizerID, ok := ctx.Value("user_id").(string)
		if !ok || organizerID == "" {
			return nil, errors.New("organizer ID not found in context")
		}

		eventDetail, err := srv.GetEventDetails(ctx, event, organizerID)
		fmt.Println("event details are",eventDetail)

		if err != nil {
			return nil, err
		}

		return eventDetail, nil
	}
}
