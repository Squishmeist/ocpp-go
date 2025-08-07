-- name: InsertChargepoint :one
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

-- name: UpdateChargepointLastHeartbeat :one
UPDATE chargepoint 
SET last_heartbeat = ?
WHERE serial_number = ?
RETURNING serial_number;
