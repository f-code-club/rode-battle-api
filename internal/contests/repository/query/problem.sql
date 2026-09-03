-- name: GetProblemsByContest :many
SELECT 
    id,
    position,
    name
FROM problems
WHERE contest_id = @contest_id
ORDER BY position ASC;

-- name: AssignProblemsToContest :execrows
UPDATE problems
SET contest_id = @contest_id
WHERE id = ANY(@problem_ids::uuid[])
  AND contest_id IS NULL;

-- name: GetProblemsForContestAssignment :many
SELECT
    id,
    contest_id
FROM problems
WHERE id = ANY(@problem_ids::uuid[]);