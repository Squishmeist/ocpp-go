-- name: InsertChargePoint :one
INSERT INTO chargepoint (
    serial_number,
    model,
    vendor,
    firmware_version,
    iicid,
    imsi,
    meter_serial_number,
    meter_type,
    last_boot,
    last_heartbeat,
    last_connected
) VALUES (?,?,?,?,?,?,?,?,?,?,?)
RETURNING *;