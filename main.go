package main

import (
	"context"
	"github.com/google/go-github/v61/github"
	"github.com/joho/godotenv"
	"github.com/permalik/utility/lo"
	"github.com/permalik/utility/repo"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

func main() {

	// TODO: drop in viper for godotenv
	lo.G(0, "launch: godotenv", nil)
	err := godotenv.Load()
	if err != nil {
		lo.G(1, "load: .env", err)
	}

	lo.G(0, "launch: go-github", nil)
	ghPAT := os.Getenv("GITHUB_PAT")
	ghClient := github.NewClient(nil).WithAuthToken(ghPAT)

	lo.G(0, "launch: redis", nil)
	redisURI := os.Getenv("REDIS_URI")
	opt, _ := redis.ParseURL(redisURI)
	rClient := redis.NewClient(opt)

	cfg := repo.Config{
		Name: "permalik",
		Org:  false,
		Ctx:  ctx,
		GC:   ghClient,
	}
	permalikRepos := repo.GHRepos(cfg)
	var ghRepos []repo.Repo
	if len(permalikRepos) > 0 {
		ghRepos = append(ghRepos, permalikRepos...)
	}

	cfg.RC = rClient
	redisKeys := repo.RedisKeys(cfg)
	if redisKeys == nil {
		for _, v := range ghRepos {
			err := repo.RedisSet(v, cfg)
			if err != nil {
				lo.G(1, "RedisSet", err)
			}
		}
		lo.G(0, "task complete", nil)
	} else {
		for _, v := range redisKeys {
			var r repo.Repo
			err := repo.RedisDelete(r, v, cfg)
			if err != nil {
				lo.G(1, "RedisRemoveOne", err)
			}
		}
		for _, v := range ghRepos {
			err := repo.RedisSet(v, cfg)
			if err != nil {
				lo.G(1, "RedisAddOne", err)
			}
		}
		lo.G(0, "task complete", nil)
	}
}
