package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v61/github"
	"github.com/joho/godotenv"
	"github.com/permalik/utility/repo"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func run(cfg repo.Config) {
	permalikRepos := repo.GHRepos(cfg)
	var ghRepos []repo.Repo
	if len(permalikRepos) > 0 {
		ghRepos = append(ghRepos, permalikRepos...)
	}

	redisKeys := repo.RedisKeys(cfg)
	if redisKeys == nil {
		for _, v := range ghRepos {
			err := repo.RedisSet(v, cfg)
			if err != nil {
				log.Fatalf("RedisSet:\n%v", err)
			}
		}
		log.Println("task complete")
	} else {
		for _, v := range redisKeys {
			var r repo.Repo
			err := repo.RedisDelete(r, v, cfg)
			if err != nil {
				log.Fatalf("RedisRemoveOne:\n%v", err)
			}
		}
		for _, v := range ghRepos {
			err := repo.RedisSet(v, cfg)
			if err != nil {
				log.Fatalf("RedisAddOne:\n%v", err)
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

	cfg := repo.Config{
		Name: "permalik",
		Org:  false,
		Ctx:  ctx,
		GC:   ghClient,
		RC:   rClient,
	}

	run(cfg)
}
