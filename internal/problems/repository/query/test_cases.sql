-- name: CreateTestCase :one
INSERT INTO test_cases (problem_id, input)
VALUES (@problem_id, @input)
RETURNING id;
