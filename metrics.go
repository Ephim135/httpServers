package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Hits int32
}

func (cfg *apiConfig) fileServerHits(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Hits: cfg.fileserverHits.Load(),
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles("admin.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
