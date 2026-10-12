-- name: GetSubmissionsByContest :many
SELECT
    s.id,
    p.name AS problem_name,
    a.name AS owner_name,
    s.language,
    s.verdict,
    s.score,
    s.created_at,
    COUNT(*) OVER() AS total_count
FROM submissions s
JOIN problems p ON p.id = s.problem_id
JOIN accounts a ON a.id = s.account_id
WHERE p.contest_id = @contest_id
  AND (@problem_id::uuid = '00000000-0000-0000-0000-000000000000'::uuid
       OR s.problem_id = @problem_id::uuid)
  AND (@account_id::uuid = '00000000-0000-0000-0000-000000000000'::uuid
       OR s.account_id = @account_id::uuid)
  AND (sqlc.narg('language')::language IS NULL OR s.language = sqlc.narg('language'))
  AND (sqlc.narg('verdict')::verdict IS NULL OR s.verdict = sqlc.narg('verdict'))
  AND (
    sqlc.narg('search')::text IS NULL
    OR p.name ILIKE '%' || sqlc.narg('search') || '%'
    OR a.name ILIKE '%' || sqlc.narg('search') || '%'
  )
ORDER BY s.created_at DESC, s.id DESC
LIMIT @page_limit OFFSET @page_offset;