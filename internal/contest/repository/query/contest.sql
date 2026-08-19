-- name: GetContestSubmissionsForRanking :many
SELECT
    a.id AS account_id,
    a.name AS account_name,
    p.id AS problem_id,
    p.position AS problem_position,
    s.language,
    s.verdict,
    s.created_at
FROM submissions s
INNER JOIN problems p ON p.id = s.problem_id
INNER JOIN contests c ON c.id = p.contest_id
INNER JOIN accounts a ON a.id = s.account_id
WHERE p.contest_id = @contest_id
  AND a.role = 'participant'
  AND a.is_banned = false
  AND s.created_at BETWEEN c.start_time AND c.end_time
ORDER BY a.id, p.position, s.created_at;

-- name: GetContestTimeRange :one
SELECT
    start_time,
    end_time
FROM contests
WHERE id = @contest_id;