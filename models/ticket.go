package models

import "time"

type Ticket struct {
	Id            int
	SlotId        int
	VehicleNumber string
	VehicleType   string
	EntryTime     time.Time
	ExitTime      *time.Time
	TotalCost     int
	Status        string // ACTIVE / COMPLETED
}

type ParkRequest struct {
	VehicleType   string `json:"vehicle_type"`
	VehicleNumber string `json:"vehicle_number"`
}


