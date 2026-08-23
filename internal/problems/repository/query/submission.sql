-- name: CreateSubmission :one
INSERT INTO submissions (problem_id, account_id, language, code)
VALUES (@problem_id, @account_id, @language, @code)
RETURNING id;

-- name: GetSubmitHistory :many
SELECT s.id, s.language, s.code, s.verdict, s.score, s.created_at
FROM submissions s
WHERE problem_id = @problem_id AND account_id = @account_id;
