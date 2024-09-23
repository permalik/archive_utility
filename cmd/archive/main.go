package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/permalik/utility/db"
	"github.com/permalik/utility/github"
	"github.com/permalik/utility/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}

	pool := db.InitDB()
	defer pool.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db.Ping(ctx)

	dropRepos(pool, ctx)
	createRepos(pool, ctx)

	allRepos := github.GithubClient()

	for _, v := range allRepos {
		insertRepos(pool, ctx, v)
	}
	selectRepos(pool, ctx)
}

func dropRepos(pool *sql.DB, ctx context.Context) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := pool.ExecContext(ctx, "DROP TABLE repos;")
	if err != nil {
		log.Fatal("unable to drop table", err)
	}
}

func createRepos(pool *sql.DB, ctx context.Context) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createQuery := `CREATE TABLE repos (
        id SERIAL PRIMARY KEY,
        owner VARCHAR(100),
        name VARCHAR(100),
        category VARCHAR(100),
        description VARCHAR(200),
        html_url VARCHAR(100),
        homepage VARCHAR(100),
        topics TEXT,
        created_at VARCHAR(10),
        updated_at VARCHAR(10),
        uid INT
    )`

	_, err := pool.ExecContext(ctx, createQuery)
	if err != nil {
		log.Fatal("unable to create table", err)
	}
}

func insertRepos(pool *sql.DB, ctx context.Context, r models.Repo) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ownerBefore, nameAfter, _ := strings.Cut(r.Data.FullName, "/")
	owner := ownerBefore
	name := nameAfter

	categoryBefore, descriptionAfter, _ := strings.Cut(r.Data.Description, ":")
	category := categoryBefore
	description := descriptionAfter

	var topics string
	for _, v := range r.Data.Topics {
		if len(topics) < 1 {
			topics = v
		} else {
			topics = fmt.Sprintf("%s,%s", topics, v)
		}
	}

	createdAt := r.Data.CreatedAt.Format("2006-01-02")
	updatedAt := r.Data.UpdatedAt.Format("2006-01-02")

	query := `
    INSERT INTO repos (
        owner,
        name,
        category,
        description,
        html_url,
        homepage,
        topics,
        created_at,
        updated_at,
        uid
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id;
    `

	result, err := pool.ExecContext(ctx, query,
		owner,
		name,
		category,
		description,
		r.Data.HTMLURL,
		r.Data.Homepage,
		topics,
		createdAt,
		updatedAt,
		r.Data.ID)
	if err != nil {
		log.Fatal("failed executing insert", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal("failed writing to db", err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d rows", rows)
	}
}

func selectRepos(pool *sql.DB, ctx context.Context) ([]byte, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := pool.QueryContext(ctx, "select * from repos;")
	if err != nil {
		log.Fatal("unable to execute select all", err)
	}
	defer rows.Close()

	var repos []models.JsonRepo
	for rows.Next() {
		var (
			id          int
			owner       string
			name        string
			category    string
			description string
			htmlURL     string
			homepage    string
			topics      string
			createdAt   string
			updatedAt   string
			uid         int
		)
		if err := rows.Scan(
			&id,
			&owner,
			&name,
			&category,
			&description,
			&htmlURL,
			&homepage,
			&topics,
			&createdAt,
			&updatedAt,
			&uid); err != nil {
			log.Fatal(err)
		}
		repo := models.JsonRepo{
			Owner:       owner,
			Name:        name,
			Category:    category,
			Description: description,
			HTMLURL:     htmlURL,
			Homepage:    homepage,
			Topics:      topics,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			UID:         uid,
		}
		repos = append(repos, repo)
	}

	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(repos)
	if err != nil {
		log.Fatalf("error marshaling to json:\n%v", err)
	}
	fmt.Println(string(jsonData))
	return jsonData, nil
}
