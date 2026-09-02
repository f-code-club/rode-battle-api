-- name: GetProblem :one
SELECT p.position, p.name, p.content, p.time_limit, p.memory_limit
FROM problems p
WHERE id = @id;

-- name: GetProblemLanguages :many
SELECT pl.language
FROM problem_languages pl
WHERE problem_id = @problem_id;

-- name: CreateProblem :one
INSERT INTO problems (name, content, checker_language, checker_path, time_limit, memory_limit)
VALUES (@name, @content, @checker_language, @checker_path, @time_limit, @memory_limit)
RETURNING id;

-- name: CreateProblemLanguage :exec
INSERT INTO problem_languages (problem_id, language)
SELECT @problem_id, unnest(@language::text[])::language;
