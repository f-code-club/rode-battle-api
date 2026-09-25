-- name: GetProblem :one
SELECT p.position, p.name, p.content, p.color_code, p.time_limit, p.memory_limit
FROM problems p
WHERE id = @id;

-- name: GetProblemLanguages :many
SELECT pl.language
FROM problem_languages pl
WHERE problem_id = @problem_id;

-- name: CreateProblem :one
INSERT INTO problems (name, content, checker_language, checker_path, time_limit, memory_limit, color_code)
VALUES (@name, @content, @checker_language, @checker_path, @time_limit, @memory_limit, @color_code)
RETURNING id;

-- name: CreateProblemLanguage :exec
INSERT INTO problem_languages (problem_id, language)
SELECT @problem_id, unnest(@language::text[])::language;

-- name: GetProblems :many
SELECT
    p.id, p.position, p.name, p.content, p.time_limit, p.memory_limit, p.color_code,
    COALESCE(array_agg(pl.language ORDER BY pl.language) FILTER (WHERE pl.language IS NOT NULL), '{}')::text[] AS languages
FROM problems p
LEFT JOIN problem_languages pl ON pl.problem_id = p.id
GROUP BY p.id, p.position, p.name, p.content, p.time_limit, p.memory_limit, p.color_code
ORDER BY p.position NULLS LAST;