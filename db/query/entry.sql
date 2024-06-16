
-- name: createEntry :one 
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: getEntry :one

SELECT * FROM entries WHERE id = $1 LIMIT 1;