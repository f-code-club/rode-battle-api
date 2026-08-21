-- name: CreateSubmission :one
INSERT INTO submissions (problem_id, account_id, language, code)
VALUES (@problem_id, @account_id, @language, @code)
RETURNING id;
