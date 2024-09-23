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

	"github.com/google/go-github/v61/github"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/permalik/utility/models"
)

var ghCtx = context.Background()

var pool *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}

	ghPAT := os.Getenv("GH_PAT")
	gc := github.NewClient(nil).WithAuthToken(ghPAT)

	permalikRepos := ghRepos(gc, "permalik", false)
	var allRepos []models.Repo
	if len(permalikRepos) > 0 {
		allRepos = append(allRepos, permalikRepos...)
	}

	dsn := os.Getenv("DSN")
	pool, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("unable to use dsn", err)
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	dbCtx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	ping(dbCtx)

	dropRepos(dbCtx)
	createRepos(dbCtx)
	for _, v := range allRepos {
		insertRepos(dbCtx, v)
	}
	selectRepos(dbCtx)
}

func ping(dbCtx context.Context) {
	dbCtx, cancel := context.WithTimeout(dbCtx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(dbCtx); err != nil {
		log.Fatalf("unable to connect to database:\n%v", err)
	}
}

func createRepos(dbCtx context.Context) {
	dbCtx, cancel := context.WithTimeout(dbCtx, 5*time.Second)
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

	_, err := pool.ExecContext(dbCtx, createQuery)
	if err != nil {
		log.Fatal("unable to create table", err)
	}
}

func dropRepos(dbCtx context.Context) {
	dbCtx, cancel := context.WithTimeout(dbCtx, 5*time.Second)
	defer cancel()

	_, err := pool.ExecContext(dbCtx, "DROP TABLE repos;")
	if err != nil {
		log.Fatal("unable to drop table", err)
	}
}

func insertRepos(dbCtx context.Context, r models.Repo) {
	dbCtx, cancel := context.WithTimeout(dbCtx, 5*time.Second)
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

	result, err := pool.ExecContext(dbCtx, query,
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

func selectRepos(dbCtx context.Context) ([]byte, error) {
	dbCtx, cancel := context.WithTimeout(dbCtx, 5*time.Second)
	defer cancel()

	rows, err := pool.QueryContext(dbCtx, "select * from repos;")
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

func parseGH(repo models.Repo, arr []models.Repo, ghData []*github.Repository) []models.Repo {
	for _, v := range ghData {
		timestampCA := v.GetCreatedAt()
		pointerCA := timestampCA.GetTime()
		createdAt := *pointerCA
		timestampUA := v.GetUpdatedAt()
		pointerUA := timestampUA.GetTime()
		updatedAt := *pointerUA
		d := models.RepoData{
			ID:          v.GetID(),
			FullName:    v.GetFullName(),
			Description: v.GetDescription(),
			HTMLURL:     v.GetHTMLURL(),
			Homepage:    v.GetHomepage(),
			Topics:      v.Topics,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		repo.Name = v.GetName()
		repo.Data = d
		arr = append(arr, repo)
	}
	return arr
}

func ghRepos(gc *github.Client, name string, isOrg bool) []models.Repo {
	var r models.Repo
	var arr []models.Repo
	listOpt := github.ListOptions{Page: 1, PerPage: 25}

	if isOrg {
		opts := &github.RepositoryListByOrgOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := gc.Repositories.ListByOrg(ghCtx, name, opts)
		if err != nil {
			log.Fatalf("github: ListByOrg\n%v", err)
		}
		if len(data) <= 0 {
			log.Fatalf("github: no data returned from GithubAll")
		}
		arr = parseGH(r, arr, data)
		return arr
	} else {
		opts := &github.RepositoryListByUserOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := gc.Repositories.ListByUser(ghCtx, name, opts)
		if err != nil {
			log.Fatalf("github: ListByUser\n%v", err)
		}
		if len(data) <= 0 {
			log.Fatalf("github: no data returned from GithubAll\n%s", name)
		}
		arr = parseGH(r, arr, data)
		return arr
	}
}
