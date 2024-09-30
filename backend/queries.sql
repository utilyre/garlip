-- name: CreateAccount :exec
INSERT INTO accounts
(created_at, updated_at, username, password, fullname)
VALUES (NOW(), NOW(), $1, $2, $3);

-- name: GetAnswersCount :many
SELECT q.stem stem, o.description description, o.correct correct, COUNT(*) num_answer
FROM optional_answers oa
JOIN options o
ON o.id = oa.option_id
JOIN questions q
ON q.id = o.question_id
WHERE q.id = $1
GROUP BY oa.option_id, q.stem, o.description, o.correct;
