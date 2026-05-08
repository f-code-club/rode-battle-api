-- name: CreateAccount :one
INSERT INTO accounts (email, password, name)
VALUES (@email, @password, @name)
RETURNING id;

-- name: GetAccountByEmail :one
SELECT id, password
FROM accounts
WHERE email = @email;
