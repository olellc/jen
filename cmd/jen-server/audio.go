package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

// GET /audio/{id}/{friendly_name}
func (app *App) audio(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	audioHandle, present := app.audios.Get(id)
	if !present {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return nil
	}

	friendly_name := chi.URLParam(r, "friendly_name")
	if friendly_name != audioHandle.FriendlyName {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", audioHandle.MediaType)
	http.ServeFile(w, r, audioHandle.Path)

	return nil
}
