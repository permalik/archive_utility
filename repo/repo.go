package repo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-github/v61/github"
	"github.com/permalik/utility/lo"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Name string
	Org  bool
	Ctx  context.Context
	GC   *github.Client
	RC   *redis.Client
}

type Data struct {
	ID          int64
	FullName    string
	Description string
	HTMLURL     string
	Homepage    string
	Topics      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repo struct {
	Name string
	Data Data
}

func parseGH(repo Repo, arr []Repo, data []*github.Repository) []Repo {

	for _, v := range data {
		timestampCA := v.GetCreatedAt()
		pointerCA := timestampCA.GetTime()
		createdAt := *pointerCA
		timestampUA := v.GetUpdatedAt()
		pointerUA := timestampUA.GetTime()
		updatedAt := *pointerUA

		d := Data{
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

func GHRepos(cfg Config) []Repo {

	var r Repo
	var arr []Repo
	listOpt := github.ListOptions{Page: 1, PerPage: 25}

	if cfg.Org {
		opt := &github.RepositoryListByOrgOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := cfg.GC.Repositories.ListByOrg(cfg.Ctx, cfg.Name, opt)
		if err != nil {
			lo.G(1, "github: ListByOrg", err)
		}
		if len(data) <= 0 {
			lo.G(0, "github: no data returned from GithubAll", cfg.Name)
			return arr
		}
		arr = parseGH(r, arr, data)
		return arr
	} else {
		opt := &github.RepositoryListByUserOptions{Type: "public", Sort: "created", ListOptions: listOpt}
		data, _, err := cfg.GC.Repositories.ListByUser(cfg.Ctx, cfg.Name, opt)
		if err != nil {
			lo.G(0, "github: ListByUser", err)
		}
		if len(data) <= 0 {
			lo.G(0, "github: no data returned from GithubAll", cfg.Name)
			return arr
		}
		arr = parseGH(r, arr, data)
		return arr
	}
}

func RedisKeys(cfg Config) []string {

	res, err := cfg.RC.Keys(cfg.Ctx, "*").Result()
	if errors.Is(err, redis.Nil) {
		lo.G(0, "RedisAll: redis.Nil: keys not found", err)
		return nil
	} else if err != nil {
		lo.G(0, "RedisAll: keys not found", err)
	}
	return res
}

func RedisSet(r Repo, cfg Config) error {

	data, err := json.Marshal(r.Data)
	if err != nil {
		lo.G(0, "RedisSet: json.Marshal", err)
		return err
	}

	err = cfg.RC.Set(cfg.Ctx, r.Name, data, 0).Err()
	if err != nil {
		lo.G(1, "RedisSet: Item not set", err)
		return err
	}
	return nil
}

func RedisDelete(r Repo, name string, cfg Config) error {

	_, err := cfg.RC.Del(cfg.Ctx, name).Result()
	if errors.Is(err, redis.Nil) {
		lo.G(0, "RedisRemoveOne: name does not exist", err)
		return nil
	}
	if err != nil {
		lo.G(0, "RedisRemoveOne: RC.Del", err)
		return err
	}
	return nil
}
