-- name: GetAccount :one
SELECT email, name, role
FROM accounts
WHERE id = @id AND is_banned = false;

-- name: GetAccountByEmail :one
SELECT id, email, name, password, role, is_banned
FROM accounts
WHERE email = @email;
