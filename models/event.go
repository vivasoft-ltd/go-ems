package models

import "time"

type Event struct {
	ID            int        `json:"id" gorm:"column:id"`
	Title         string     `json:"title" gorm:"column:title"`
	Description   *string    `json:"description" gorm:"column:description"`
	Location      *string    `json:"location" gorm:"column:location"`
	StartTime     *time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime       *time.Time `json:"end_time" gorm:"column:end_time"`
	IsPublic      bool       `json:"is_public" gorm:"column:is_public"`
	AttendeeLimit int        `json:"attendee_limit" gorm:"column:attendee_limit"`
	CreatedBy     int        `json:"created_by" gorm:"column:created_by"`
	CreatedAt     time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"-" gorm:"column:updated_at"`
	Attendees     []User     `json:"attendees" gorm:"many2many:event_attendees;"`
}

type EventAttendee struct {
	EventID  int `json:"-" gorm:"column:event_id"`
	UserID   int `json:"-" gorm:"column:user_id"`
	StatusID int `json:"-" gorm:"column:status_id"`
}
