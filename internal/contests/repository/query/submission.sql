-- name: GetSubmissionsByContest :many
SELECT
    s.id,
    p.name AS problem_name,
    a.name AS owner_name,
    s.language,
    s.verdict,
    s.score,
    s.created_at
FROM submissions s
JOIN problems p ON p.id = s.problem_id
JOIN accounts a ON a.id = s.account_id
WHERE p.contest_id = @contest_id
ORDER BY s.created_at DESC;