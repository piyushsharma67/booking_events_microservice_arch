package repository

import (
	"context"

	"github.com/piyushsharma67/events_booking/services/events_service/database"
	"github.com/piyushsharma67/events_booking/services/events_service/models"
)

type EventRepository struct {
	db database.Database
}

func NewRepos(db database.Database)*EventRepository{
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GenerateEvent(ctx context.Context, event *models.EventDocument) (*models.EventDocument, error) {

	return r.db.GenerateEvent(ctx, event)
}

func (r *EventRepository) DeleteEventByOrganiser(ctx context.Context, id any) error {
	return r.db.DeleteEvent(ctx, 1)
}

func (r *EventRepository) DeleteEventByAdmin(ctx context.Context, id any) error {
	return r.db.DeleteEvent(ctx, 1)
}

func (r *EventRepository) UpdateEventByOrganiser(ctx context.Context, event *models.EventDocument) (*models.EventDocument, error) {
	return r.db.UpdateEvent(ctx, event)
}

func (r *EventRepository) UpdateEventByAdmin(ctx context.Context, event *models.EventDocument) (*models.EventDocument, error) {
	return r.db.UpdateEvent(ctx, event)
}

func (r *EventRepository) GetEvent(ctx context.Context, id interface{}) (*models.EventDocument, error) {
	return r.db.GetEvent(ctx, id)
}
