package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"parking_slot/services"
)

type handlers struct {
	service services.Service
}

func NewHandler(s services.Service) *handlers {
	return &handlers{service: s}
}

type parkRequest struct {
	VehicleType   string `json:"vehicle_type"`
	VehicleNumber string `json:"vehicle_number"`
}

type unparkRequest struct {
	TicketId int `json:"ticket_id"`
}

func (h *handlers) Park(w http.ResponseWriter, r *http.Request) {
	var req parkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := h.service.Park(r.Context(), req.VehicleNumber, req.VehicleType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func (h *handlers) UnPark(w http.ResponseWriter, r *http.Request) {
	var req unparkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	unparked, err := h.service.Unpark(r.Context(), req.TicketId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"ticket_id":%d,"unparked":%t}`, req.TicketId, unparked)))
}

func (h *handlers) Available(w http.ResponseWriter, r *http.Request) {
	vehicleType := r.URL.Query().Get("vehicle_type")

	slots, err := h.service.GetAvailableSlots(r.Context(), vehicleType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}
