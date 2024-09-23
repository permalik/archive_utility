package main

import "net/http"

func (app *application) GetRepos(w http.ResponseWriter, r *http.Request) {
	jsonData, err := RepoService(app.pool, app.ctx)
	if err != nil {
		http.Error(w, "error fetching repos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
