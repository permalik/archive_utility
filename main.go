package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type config struct {
	Name string
	Org  bool
	Ctx  context.Context
	GC   *github.Client
	RC   *redis.Client
}

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

func GHRepos(cfg config) []repo {
	var r repo
	var arr []repo
	listOpt := github.ListOptions{Page: 1, PerPage: 25}

	if cfg.Org {
		opt := &github.RepositoryListByOrgOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := cfg.GC.Repositories.ListByOrg(cfg.Ctx, cfg.Name, opt)
		if err != nil {
			log.Fatalf("github: ListByOrg\n%v", err)
		}
		if len(data) <= 0 {
			log.Fatalf("github: no data returned from GithubAll")
		}
		arr = parseGH(r, arr, data)
		return arr
	} else {
		opt := &github.RepositoryListByUserOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := cfg.GC.Repositories.ListByUser(cfg.Ctx, cfg.Name, opt)
		if err != nil {
			log.Fatalf("github: ListByUser\n%v", err)
		}
		if len(data) <= 0 {
			log.Fatalf("github: no data returned from GithubAll\n%s", cfg.Name)
		}
		arr = parseGH(r, arr, data)
		return arr
	}
}

func run(cfg config) {
	permalikRepos := GHRepos(cfg)
	var ghRepos []repo
	if len(permalikRepos) > 0 {
		ghRepos = append(ghRepos, permalikRepos...)
	}

	redisKeys, err := cfg.RC.Keys(cfg.Ctx, "*").Result()
	if errors.Is(err, redis.Nil) {
		log.Fatalf("RedisAll: redis.Nil: keys not found\n%v", err)
	} else if err != nil {
		log.Fatalf("RedisAll: keys not found\n%v", err)
	}

	if redisKeys == nil {
		for _, v := range ghRepos {
			data, err := json.Marshal(v.Data)
			if err != nil {
				log.Fatalf("RedisSet: json.Marshal\n%v", err)
			}
			err = cfg.RC.Set(cfg.Ctx, v.Name, data, 0).Err()
			if err != nil {
				log.Fatalf("RedisSet:\n%v", err)
			}
		}
		log.Println("task complete")
	} else {
		for _, v := range redisKeys {
			_, err := cfg.RC.Del(cfg.Ctx, v).Result()
			if errors.Is(err, redis.Nil) {
				log.Fatalf("RedisRemoveOne: name does not exist\n%v", err)
			}
			if err != nil {
				log.Fatalf("RedisRemoveOne: RC.Del\n%v", err)
			}
			if err != nil {
				log.Fatalf("RedisRemoveOne:\n%v", err)
			}
		}
		for _, v := range ghRepos {
			data, err := json.Marshal(v.Data)
			if err != nil {
				log.Fatalf("RedisSet: json.Marshal\n%v", err)
			}
			err = cfg.RC.Set(cfg.Ctx, v.Name, data, 0).Err()
			if err != nil {
				log.Fatalf("RedisSet:\n%v", err)
			}
		}
		log.Println("task complete")
	}
}

func main() {

	log.Println("launch: godotenv")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load .env:\n%v", err)
	}

	log.Println("launch: go-github")
	ghPAT := os.Getenv("GITHUB_PAT")
	ghClient := github.NewClient(nil).WithAuthToken(ghPAT)

	log.Println("launch: redis")
	redisURI := os.Getenv("REDIS_URI")
	opt, _ := redis.ParseURL(redisURI)
	rClient := redis.NewClient(opt)

	cfg := config{
		Name: "permalik",
		Org:  false,
		Ctx:  ctx,
		GC:   ghClient,
		RC:   rClient,
	}

	run(cfg)
}
