-- name: GetProblem :one
SELECT p.position, p.name, p.content_path, p.time_limit, p.memory_limit
FROM problems p
WHERE id = @id;

-- name: GetProblemLanguages :many
SELECT pl.language
FROM problem_languages pl
WHERE problem_id = @problem_id;