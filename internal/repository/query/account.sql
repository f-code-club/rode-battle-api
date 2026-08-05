-- name: CreateAccount :one
INSERT INTO accounts (email, name, password, role)
VALUES (@email, @name, @password, @role)
RETURNING id;

-- name: GetAccountByEmail :one
SELECT id, email, name, password, role, is_banned
FROM accounts
WHERE email = @email;
