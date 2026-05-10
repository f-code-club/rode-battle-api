-- name: CreateAccount :one
INSERT INTO accounts (email, password, name, school, student_id, phone_number)
VALUES (@email, @password, @name, @school, @student_id, @phone_number)
RETURNING id;

-- name: IsEmailRegistered :one
SELECT EXISTS(
    SELECT 1
    FROM accounts
    WHERE email = @email
);