-- name: GetProblemsByContest :many
SELECT 
    id,
    position,
    name
FROM problems
WHERE contest_id = @contest_id
ORDER BY position ASC;

-- name: AssignProblemToContest :exec
UPDATE problems
SET contest_id = @contest_id, position = @position
WHERE id = @id;
