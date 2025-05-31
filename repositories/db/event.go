package db

import (
	"errors"
	"fmt"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Repository) CreateEvent(event *models.Event) (*models.Event, error) {
	qry := repo.client.Create(event)
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error creating event: %w", qry.Error))
		return nil, qry.Error
	}

	return event, nil
}

func (repo *Repository) ListEvents(filter *types.EventFilter, limit, offset int) ([]*models.Event, int, error) {
	var events []*models.Event
	var count int64
	query := repo.client.Model(&models.Event{})
	repo.applyEventFilter(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	result := query.Offset(offset).Limit(limit).Find(&events)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if result.RowsAffected == 0 {
		logger.Error("no events found")
		return nil, 0, errutil.ErrRecordNotFound
	}

	return events, int(count), nil
}
func (repo *Repository) applyEventFilter(query *gorm.DB, filter *types.EventFilter) {
	if filter != nil {
		if filter.IsPublic != nil && *filter.IsPublic {
			query = query.Where("is_public = ? and end_time > now()", true)
		}
		if filter.CreatedBy != nil {
			query = query.Where("created_by = ?", filter.CreatedBy)
		}
		if filter.Attendee != nil {
			query = query.Where("(is_public = ? and end_time > now() ) or (id in (select event_id from event_attendees where user_id = ?))", true, filter.Attendee)
		}
	}
}

func (repo *Repository) ReadEventByID(id int) (*models.Event, error) {
	var event models.Event
	qry := repo.client.First(&event, id)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Errorf("event with ID %d not found", id))
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error getting event by ID: %w", qry.Error))
		return nil, qry.Error
	}

	return &event, nil
}

func (repo *Repository) UpdateEvent(event *models.Event) (*models.Event, error) {
	qry := repo.client.Where("id = ?", event.ID).Updates(event)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Errorf("no event found with ID %d", event.ID))
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error updating event: %w", qry.Error))
		return nil, qry.Error
	}
	return event, nil
}

func (repo *Repository) DeleteEvent(id int) error {
	qry := repo.client.Where("id = ?", id).Delete(&models.Event{})
	if qry.RowsAffected == 0 {
		logger.Error(fmt.Errorf("no event found with ID %d", id))
		return errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error deleting event: %w", qry.Error))
		return qry.Error
	}
	return nil
}
func (repo *Repository) ReadInvitationByEventAndUser(eventID, userID int) (*models.EventAttendee, error) {
	var invitation models.EventAttendee
	qry := repo.client.Where("event_id = ? AND user_id = ?", eventID, userID).First(&invitation)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Errorf("no invitation found with event ID %d and user ID %d", eventID, userID))
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error getting invitation: %w", qry.Error))
		return nil, qry.Error
	}
	return &invitation, nil
}
func (repo *Repository) UpsertInvitation(invitation *models.EventAttendee) error {
	qry := repo.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "event_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status_id"}),
	}).Create(invitation)
	if qry.Error != nil {
		logger.Error(fmt.Errorf("error upserting invitation: %w", qry.Error))
		return qry.Error
	}
	return nil
}
func (repo *Repository) GetEventAttendeesCount(eventID int) (int, error) {
	var count int64
	if err := repo.client.Model(&models.EventAttendee{}).Where("event_id = ? and status_id = 2", eventID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
