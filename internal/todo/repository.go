package todo

import (
	"context"
	"database/sql"
)

type Repository struct {
	db Execer
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewRepository(connection *sql.DB) *Repository {
	return &Repository{db: connection}
}

func NewRepositoryTx(connection *sql.Tx) *Repository {
	return &Repository{db: connection}
}

func (r *Repository) Create(ctx context.Context, todo Todo) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO todos(title, description, done) VALUES($1, $2, $3)", todo.Title, todo.Description, todo.Done)
	return err
}

func (r *Repository) GetAll(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT (id, title, description, done) FROM todos")
	defer rows.Close()
	if err != nil {
		return []Todo{}, err
	}
	var results []Todo
	for rows.Next() {
		var result Todo
		err = rows.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
		if err != nil {
			break
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, err
}

func (r *Repository) GetById(ctx context.Context, id int) (Todo, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, title, description, done FROM todos WHERE id=$1", id)

	var result Todo
	err := row.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
	if err != nil {
		return Todo{}, err
	}
	return result, nil
}

func (r *Repository) Update(ctx context.Context, todo Todo) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE todos SET title=$1, description=$2, done=$3 WHERE id=$4", todo.Title, todo.Description, todo.Done, todo.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx,
		"DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
