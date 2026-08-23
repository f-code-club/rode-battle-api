-- name: GetProblem :one
SELECT p.position, p.name, p.content_path, p.time_limit, p.memory_limit
FROM problems p
WHERE id = @id;

-- name: GetProblemLanguages :many
SELECT pl.language
FROM problem_languages pl
WHERE problem_id = @problem_id;

-- name: GetProblemsByContest :many
SELECT 
    id,
    position,
    name
FROM problems
WHERE contest_id = @contest_id
ORDER BY position ASC;