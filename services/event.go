package services

import (
	"errors"
	"github.com/vivasoft-ltd/go-ems/consts"

	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
)

type EventServiceImpl struct {
	eventRepo domain.EventRepository
	userRepo  domain.UserRepository
}

func NewEventServiceImpl(eventRepo domain.EventRepository, userRepo domain.UserRepository) *EventServiceImpl {
	return &EventServiceImpl{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (svc *EventServiceImpl) CreateEvent(eventReq *types.CreateEventRequest) (*types.CreateEventResponse, error) {
	event := eventReq.ToEvent()
	if !eventReq.IsPublic {
		users, err := svc.userRepo.ReadUsersByIDs(eventReq.Attendees)
		if err != nil {
			return nil, err
		}
		event.Attendees = users
	}
	createdEvent, err := svc.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &types.CreateEventResponse{
		Message: "Event created",
		Event:   createdEvent,
	}, nil
}

func (svc *EventServiceImpl) ListEvents(req types.ListEventRequest, user *types.CurrentUser) (*types.PaginatedEventResponse, error) {
	offset := (req.Page - 1) * req.Limit
	filter := svc.getEventFilter(user)
	events, count, err := svc.eventRepo.ListEvents(filter, req.Limit, offset)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return &types.PaginatedEventResponse{}, nil
	}
	if err != nil {
		return nil, err
	}
	response := &types.PaginatedEventResponse{
		Page:   req.Page,
		Limit:  req.Limit,
		Total:  count,
		Events: events,
	}
	return response, nil
}
func (svc *EventServiceImpl) getEventFilter(user *types.CurrentUser) *types.EventFilter {
	filter := &types.EventFilter{}
	if user == nil {
		t := true
		filter.IsPublic = &t
		return filter
	}
	if user.HasPermission(consts.PermissionFetchAllEvent) {
		return filter
	}
	if user.HasPermission(consts.PermissionFetchOwnEvent) {
		filter.CreatedBy = &user.ID
	}
	if user.HasPermission(consts.PermissionFetchInvitedEvent) {
		filter.Attendee = &user.ID
	}
	return filter

}

func (svc *EventServiceImpl) ReadEventByID(id int) (*models.Event, error) {
	event, err := svc.eventRepo.ReadEventByID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (svc *EventServiceImpl) UpdateEvent(eventReq *types.UpdateEventRequest) (*types.UpdateEventResponse, error) {
	existingEvent, err := svc.eventRepo.ReadEventByID(eventReq.ID)
	if err != nil {
		return nil, err
	}
	if existingEvent == nil {
		return nil, errutil.ErrRecordNotFound
	}

	event := eventReq.ToEvent()
	updatedEvent, err := svc.eventRepo.UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return &types.UpdateEventResponse{
		Message: "Event updated",
		Event:   updatedEvent,
	}, nil
}

func (svc *EventServiceImpl) DeleteEvent(id int) (*types.DeleteEventResponse, error) {
	err := svc.eventRepo.DeleteEvent(id)
	if err != nil {
		return nil, err
	}
	return &types.DeleteEventResponse{
		Message: "Event deleted",
	}, nil
}
func (svc *EventServiceImpl) RsvpEvent(req types.RsvpRequest) error {
	event, err := svc.eventRepo.ReadEventByID(req.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errutil.ErrRecordNotFound
	}
	if !event.IsPublic {
		invitation, err := svc.eventRepo.ReadInvitationByEventAndUser(req.EventID, req.UserID)
		if err != nil {
			return err
		}
		if invitation == nil {
			return errutil.ErrRecordNotFound
		}
		invitation.StatusID = req.StatusID
		if err := svc.eventRepo.UpsertInvitation(invitation); err != nil {
			return err
		}
		return nil

	}
	count, err := svc.eventRepo.GetEventAttendeesCount(req.EventID)
	if err != nil {
		return err
	}
	if count >= event.AttendeeLimit {
		return errutil.ErrEventFull
	}

	newInvitation := models.EventAttendee{
		EventID:  req.EventID,
		UserID:   req.UserID,
		StatusID: req.StatusID,
	}

	if err := svc.eventRepo.UpsertInvitation(&newInvitation); err != nil {
		return err
	}
	return nil
}
