package todo

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const testDsn = "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable"

// TODO: rewrite tests so they don't depend on db data
func TestRepository_Create(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = testDsn
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	repo := NewRepositoryTx(tx)

	todo := Todo{Title: "unit test", Description: "check insert", Done: false}
	err = repo.Create(ctx, todo)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_GetAll(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = testDsn
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	repo := NewRepositoryTx(tx)

	todos, err := repo.GetAll(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 3 {
		t.Fatalf("wrong length: %d", len(todos))
	}
	example := Todo{Id: 1, Title: "first todo", Description: "first todo just to fill up DB", Done: false}
	if todos[0] != example {
		t.Fatalf("wrong first record: %s", todos[0].Description)
	}
}

func TestRepository_GetById(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = testDsn
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	repo := NewRepositoryTx(tx)

	todo, err := repo.GetById(ctx, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	example := Todo{Id: 1, Title: "first todo", Description: "first todo just to fill up DB", Done: false}
	if todo != example {
		t.Fatalf("wrong first record: %s", todo.Description)
	}
}

func TestRepository_Update(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = testDsn
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	repo := NewRepositoryTx(tx)

	todo := Todo{Id: 1, Title: "unit test", Description: "check insert", Done: false}
	err = repo.Update(ctx, todo)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changedTodo, err := repo.GetById(ctx, 1)
	if changedTodo.Title != todo.Title {
		t.Fatalf("wrong insertion: %s", changedTodo.Title)
	}
}

func TestRepository_Delete(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = testDsn
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	repo := NewRepositoryTx(tx)

	err = repo.Delete(ctx, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	todo, err := repo.GetById(ctx, 1)

	if errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("unexpected error: %v", err)
	}
	example := Todo{}
	if todo != example {
		t.Fatalf("wrong first record: %s", todo.Description)
	}
}
