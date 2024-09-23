package models

import "time"

type RepoData struct {
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
	Data RepoData
}

type JsonRepo struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	HTMLURL     string `json:"htmlurl"`
	Homepage    string `json:"homepage"`
	Topics      string `json:"topics"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	UID         int    `json:"uid"`
}
