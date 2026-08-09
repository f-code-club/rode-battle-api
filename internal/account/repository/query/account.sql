-- name: CreateAccount :one
INSERT INTO accounts (email, name, password, role)
VALUES (@email, @name, @password, @role)
RETURNING id;
