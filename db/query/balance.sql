-- name: GetBalance :many
SELECT c.code , sum(tr.amount) AS balance 
FROM public.transfers AS tr
JOIN currencies c ON tr.currency_id = c.id
WHERE tr.user_id = $1
GROUP BY c.code;
