package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/piyushsharma67/events_booking/services/events_service/infra"
	"github.com/piyushsharma67/events_booking/services/events_service/models"
	"github.com/piyushsharma67/events_booking/services/events_service/que"
	"github.com/piyushsharma67/events_booking/services/events_service/repository"
	"github.com/redis/go-redis/v9"
)

type EventService struct {
	Repository repository.EventRepository
	Publisher  *que.QuePublisher
	Redis      *redis.Client
	CacheTTL   time.Duration
}

func GetEventService(repository repository.EventRepository, publisher *que.QuePublisher,redis *redis.Client) *EventService {
	return &EventService{
		Repository: repository,
		Publisher:  publisher,
		Redis: redis,
		CacheTTL: infra.DefaultTTL(),
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	req *models.CreateEventRequest,
	organiserId string,
) (*models.EventDocument, error) {

	eventDoc, err := models.MapCreateRequestToDocument(req, organiserId)
	if err != nil {
		return nil, err
	}

	created, err := s.Repository.GenerateEvent(ctx, eventDoc)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *EventService) CreateEventAndGenerateSeats(
	ctx context.Context,
	req *models.CreateEventRequest,
	organiserId string,
) (*models.EventDocument, error) {

	// 1️⃣ Create event document
	eventDoc, err := models.MapCreateRequestToDocument(req, organiserId)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Persist event
	created, err := s.Repository.GenerateEvent(ctx, eventDoc)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Marshal event for RabbitMQ
	payload, err := json.Marshal(created)
	if err != nil {
		return created, err
	}

	// 3️⃣ Publish seats generation event
	err = s.Publisher.Publish(ctx, payload)
	if err != nil {
		// IMPORTANT: do NOT rollback DB here
		// Log & retry via DLQ if needed
		return created, err
	}

	return created, nil
}
func (s *EventService) GetEventDetails(
	ctx context.Context,
	req *models.GetEventDetailtRequest,
	organiserId string,
) (*models.EventDetailResult, error) {

	cacheKey := "event:details:" + req.EventID

	// 1️⃣ Try cache
	if s.Redis != nil {
		cached, err := s.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var resp models.EventDetailResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return &models.EventDetailResult{
					Data:   &resp,
					Source: models.SourceCache,
				}, nil
			}
		}
	}

	// 2️⃣ Fetch from DB
	created, err := s.Repository.GetEvent(ctx, req.EventID)
	if err != nil {
		return nil, err
	}

	resp := &models.EventDetailResponse{
		EventBase: models.EventBase{
			Title:       created.Title,
			Description: created.Description,
			ImageURL:    created.ImageURL,
			Location:    created.Location,
			StartTime:   created.StartTime,
			EndTime:     created.EndTime,
			Rows:        models.MapSeatingRows(created.Rows),
		},
	}

	// 3️⃣ Store in cache (best-effort)
	if s.Redis != nil {
		if payload, err := json.Marshal(resp); err == nil {
			_ = s.Redis.Set(ctx, cacheKey, payload, s.CacheTTL).Err()
		}
	}

	return &models.EventDetailResult{
		Data:   resp,
		Source: models.SourceDB,
	}, nil
}
