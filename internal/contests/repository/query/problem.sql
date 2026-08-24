-- name: GetProblemsByContest :many
SELECT 
    id,
    position,
    name
FROM problems
WHERE contest_id = @contest_id
ORDER BY position ASC;