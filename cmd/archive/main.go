package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
)

var ctxBG = context.Background()

var pool *sql.DB

type data struct {
	ID          int64
	FullName    string
	Description string
	HTMLURL     string
	Homepage    string
	Topics      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type repo struct {
	Name string
	Data data
}

// TODO: get repos from gh
// TODO: store repos in pg
// TODO: send email to pm
// TODO: serve repos from api

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load .env:\n%v", err)
	}

	ghPAT := os.Getenv("GH_PAT")
	gc := github.NewClient(nil).WithAuthToken(ghPAT)

	permalikRepos := ghRepos(gc, "permalik", false)
	var allRepos []repo
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

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	ping(ctx)

	/*
		type data struct {
			ID          int64
			FullName    string
			Description string
			HTMLURL     string
			Homepage    string
			Topics      []string
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		type repo struct {
			Name string
			Data data
		}
	*/

	for _, v := range allRepos {
		fmt.Println(v.Name)
		fmt.Println(v.Data.ID)
		fmt.Println(v.Data.FullName)
		fmt.Println(v.Data.Description)
		fmt.Println(v.Data.HTMLURL)
		fmt.Println(v.Data.Homepage)
		fmt.Println(v.Data.CreatedAt)
		fmt.Println(v.Data.UpdatedAt)

		for _, v := range v.Data.Topics {
			fmt.Println(v)
		}
	}

	// dropRepos(ctx)
	createRepos(ctx)
	// selectRepos(ctx)
}

func ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database:\n%v", err)
	}
}

func createRepos(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func dropRepos(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := pool.ExecContext(ctx, "drop table repos;")
	if err != nil {
		log.Fatal("unable to drop table", err)
	}
}

func selectRepos(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := pool.QueryContext(ctx, "select * from repos;")
	if err != nil {
		log.Fatal("unable to execute select all", err)
	}
	defer rows.Close()

	repos := make([]string, 0)
	for rows.Next() {
		var id int
		var repo string
		if err := rows.Scan(&id, &repo); err != nil {
			log.Fatal(err)
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

	fmt.Printf("people: %s", strings.Join(repos, ","))
}

func parseGH(repo repo, arr []repo, ghData []*github.Repository) []repo {
	for _, v := range ghData {
		timestampCA := v.GetCreatedAt()
		pointerCA := timestampCA.GetTime()
		createdAt := *pointerCA
		timestampUA := v.GetUpdatedAt()
		pointerUA := timestampUA.GetTime()
		updatedAt := *pointerUA
		d := data{
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

func ghRepos(gc *github.Client, name string, isOrg bool) []repo {
	var r repo
	var arr []repo
	listOpt := github.ListOptions{Page: 1, PerPage: 25}

	if isOrg {
		opts := &github.RepositoryListByOrgOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := gc.Repositories.ListByOrg(ctxBG, name, opts)
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
		data, _, err := gc.Repositories.ListByUser(ctxBG, name, opts)
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
