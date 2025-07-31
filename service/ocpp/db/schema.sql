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

-- Message Table
CREATE TABLE message (
    uuid TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('REQUEST', 'CONFIRMATION')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL CHECK (action IN ('Heartbeat', 'BootNotification')),
    payload TEXT NOT NULL,
    PRIMARY KEY (uuid, type)
);