-- name: CreateAccount :one
INSERT INTO accounts (email, password, name, school, student_id, phone_number)
VALUES (@email, @password, @name, @school, @student_id, @phone_number)
RETURNING id;

-- name: GetAccountByEmail :one
SELECT id, email, name, role, password, is_verified, is_banned
FROM accounts
WHERE email = @email;