package handler

import "github.com/go-chi/chi/v5"

func Routes(r *chi.Mux) {
	r.Get("/extract-text", PythonTextExtractorHandler)
}