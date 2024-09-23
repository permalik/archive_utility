package main

import "net/http"

func (app *application) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET 127.0.0.1/repos", app.GetRepos)
	return mux
}
