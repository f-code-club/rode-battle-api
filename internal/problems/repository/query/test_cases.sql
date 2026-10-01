-- name: CreateTestCase :one
INSERT INTO test_cases (problem_id, input_path)
VALUES (@problem_id, @input_path)
RETURNING id;
