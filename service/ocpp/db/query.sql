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

-- name: InsertMessage :one
INSERT INTO message (
    uuid,
    type,
    action,
    payload
) VALUES (?,?,?,?)
RETURNING *;

-- name: GetMessagesByUuid :one
SELECT *
FROM message
WHERE uuid = ?;

-- name: GetRequestMessageByUuid :one
SELECT *
FROM message
WHERE uuid = ? AND type = 'REQUEST';