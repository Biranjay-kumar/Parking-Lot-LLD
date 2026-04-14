package repo

import (
	"context"
	"parking_slot/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Park(ctx context.Context, vehicleType string, vehicleNumber string) (models.Ticket, error)
	Unpark(ctx context.Context, ticketID int) (bool, error)
	GetAvailableSlots(ctx context.Context, vehicleType string) ([]models.Slot, error)
}

type repo struct {
	pool *pgxpool.Pool
}

func NewRepo(p *pgxpool.Pool) Repo {
	return &repo{pool: p}
}

func (r *repo) Park(ctx context.Context, vehicleType string, vehicleNumber string) (models.Ticket, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Ticket{}, err
	}
	defer tx.Rollback(ctx)

	q := `
	SELECT id, floor_no, slot_no 
	FROM parking_slots
	WHERE vehicle_type = $1 AND status = 'AVAILABLE'
	ORDER BY floor_no, slot_no
	LIMIT 1
	FOR UPDATE;
	`

	var slotID, floor, slot int
	err = tx.QueryRow(ctx, q, vehicleType).Scan(&slotID, &floor, &slot)
	if err != nil {
		return models.Ticket{}, err
	}

	q2 := `UPDATE parking_slots SET status = 'OCCUPIED' WHERE id = $1`
	_, err = tx.Exec(ctx, q2, slotID)
	if err != nil {
		return models.Ticket{}, err
	}

	q3 := `
	INSERT INTO tickets (vehicle_type, vehicle_number, parking_slot_id, entry_time, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, parking_slot_id, entry_time, status;
	`

	var ticket models.Ticket
	err = tx.QueryRow(ctx, q3, vehicleType, vehicleNumber, slotID, time.Now(), "ACTIVE").
		Scan(&ticket.Id, &ticket.SlotId, &ticket.EntryTime, &ticket.Status)
	if err != nil {
		return models.Ticket{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (r *repo) Unpark(ctx context.Context, ticketID int) (bool, error) {
	txn, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer txn.Rollback(ctx)

	// get slot id from ticket
	var slotID int
	var entryTime time.Time
	q1 := `SELECT parking_slot_id, entry_time FROM tickets WHERE id = $1 AND status = 'ACTIVE'`
	err = txn.QueryRow(ctx, q1, ticketID).Scan(&slotID, &entryTime)
	if err != nil {
		return false, err
	}

	exitTime := time.Now()
	duration := exitTime.Sub(entryTime).Hours()
	totalCost := int(duration * 50) // 50 per hour

	// update ticket
	q2 := `UPDATE tickets SET status = 'COMPLETED', exit_time = $1, total_cost = $2 WHERE id = $3`
	_, err = txn.Exec(ctx, q2, exitTime, totalCost, ticketID)
	if err != nil {
		return false, err
	}

	// free the slot
	q3 := `UPDATE parking_slots SET status = 'AVAILABLE' WHERE id = $1`
	_, err = txn.Exec(ctx, q3, slotID)
	if err != nil {
		return false, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repo) GetAvailableSlots(ctx context.Context, vehicleType string) ([]models.Slot, error) {
	q := `SELECT id, slot_no, floor_no, vehicle_type, status FROM parking_slots WHERE status = 'AVAILABLE' AND vehicle_type = $1 ORDER BY floor_no, slot_no`

	rows, err := r.pool.Query(ctx, q, vehicleType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		err = rows.Scan(&s.Id, &s.Slot, &s.Floor, &s.VehicleType, &s.Status)
		if err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
