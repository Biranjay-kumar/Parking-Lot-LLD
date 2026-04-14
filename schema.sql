CREATE TABLE parking_slots (
    id SERIAL PRIMARY KEY,
    floor_no INT,
    slot_no INT,
    vehicle_type VARCHAR(20),
    status VARCHAR(20)
);

CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    parking_slot_id INT REFERENCES parking_slots(id),
    vehicle_number VARCHAR(20),
    vehicle_type VARCHAR(20),
    entry_time TIMESTAMP,
    exit_time TIMESTAMP,
    total_cost INT,
    status VARCHAR(20)
);
