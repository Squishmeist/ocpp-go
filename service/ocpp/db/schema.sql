-- Chargepoint Table
CREATE TABLE chargepoint (
    serial_number TEXT PRIMARY KEY NOT NULL,
    model TEXT NOT NULL,
    vendor TEXT NOT NULL,
    firmware_version TEXT NOT NULL,
    iicid TEXT,
    imsi TEXT,
    meter_serial_number TEXT,
    meter_type TEXT,
    last_boot TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_heartbeat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_connected TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
