-- name: CreateAccount :exec
INSERT INTO accounts
(created_at, updated_at, username, password, fullname)
VALUES (NOW(), NOW(), $1, $2, $3);

-- name: UpdateAccount :exec
UPDATE "accounts"
SET
    "username" = COALESCE(NULLIF(@username::VARCHAR(50), ''), "username"),
    "fullname" = COALESCE(NULLIF(@fullname::VARCHAR(100), ''), "fullname"),
    "bio" = COALESCE(NULLIF(@bio::TEXT, ''), "bio")
WHERE "id" = $1;

-- name: GetAccountPasswordByUsername :one
SELECT password
FROM accounts
WHERE username = $1
LIMIT 1;

-- name: GetAnswersCount :many
SELECT q.stem stem, o.description description, o.correct correct, COUNT(*) num_answer
FROM answers oa
JOIN options o
ON o.id = oa.option_id
JOIN questions q
ON q.id = o.question_id
WHERE q.id = $1
GROUP BY oa.option_id, q.stem, o.description, o.correct;
