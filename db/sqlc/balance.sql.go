// Code generated by sqlc. DO NOT EDIT.
// source: balance.sql

package db

import (
	"context"
)

const getBalance = `-- name: GetBalance :many
SELECT c.code , sum(tr.amount) AS balance 
FROM public.transfers AS tr
JOIN currencies c ON tr.currency_id = c.id
WHERE tr.user_id = $1
GROUP BY c.code
`

type GetBalanceRow struct {
	Code    string `json:"code"`
	Balance int64  `json:"balance"`
}

func (q *Queries) GetBalance(ctx context.Context, userID int64) ([]GetBalanceRow, error) {
	rows, err := q.db.QueryContext(ctx, getBalance, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBalanceRow{}
	for rows.Next() {
		var i GetBalanceRow
		if err := rows.Scan(&i.Code, &i.Balance); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}