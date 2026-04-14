package services

import (
	"context"
	"parking_slot/models"
	"parking_slot/repo"
)

type Service interface {
	Park(ctx context.Context, vehicleNumber string, vehicleType string) (models.Ticket, error)
	Unpark(ctx context.Context, ticketId int) (bool, error)
	GetAvailableSlots(ctx context.Context, vehicleType string) ([]models.Slot, error)
}

type services struct {
	svr repo.Repo
}

func NewService(s repo.Repo) Service {
	return &services{svr: s}
}

func (s *services) Park(ctx context.Context, vehicleNumber string, vehicleType string) (models.Ticket, error) {
	return s.svr.Park(ctx, vehicleType, vehicleNumber)
}

func (s *services) Unpark(ctx context.Context, ticketId int) (bool, error) {
	return s.svr.Unpark(ctx, ticketId)
}

func (s *services) GetAvailableSlots(ctx context.Context, vehicleType string) ([]models.Slot, error) {
	return s.svr.GetAvailableSlots(ctx, vehicleType)
}
