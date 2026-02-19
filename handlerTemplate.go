package main

import (
	"net/http"
)

func (cfg *config) handlerTemplate(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Add("Content-Type", "text/html")
		err := cfg.temp.ExecuteTemplate(w, templateName, nil)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}

	}
}
