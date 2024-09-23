package github

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v61/github"
	"github.com/permalik/utility/models"
)

var ctx = context.Background()

func GithubClient() []models.Repo {
	ghPAT := os.Getenv("GH_PAT")
	gc := github.NewClient(nil).WithAuthToken(ghPAT)

	permalikRepos := ghRepos(gc, "permalik", false)
	var allRepos []models.Repo
	if len(permalikRepos) > 0 {
		allRepos = append(allRepos, permalikRepos...)
	}
	return allRepos
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
		data, _, err := gc.Repositories.ListByOrg(ctx, name, opts)
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
		data, _, err := gc.Repositories.ListByUser(ctx, name, opts)
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
