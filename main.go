package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
)

var (
	Engine *chi.Mux

	allowAPI jsoniter.API

	relaxedAPI = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  false,
	}.Froze()
)

func main() {
	allowAPI = relaxedAPI

	Engine = chi.NewRouter()
	Engine.Use(middleware.Logger)
	Engine.Use(middleware.Recoverer)

	Registration(Engine)

	fmt.Println("we are listening live at http://localhost:8080")
	if err := http.ListenAndServe(":8080", Engine); err != nil {
		log.Fatal(err)
	}
}

func Registration(router *chi.Mux) {
	router.Route("/org", func(r chi.Router) {
		r.Post("/", createOrg)

		r.Route("/{orgId}", func(r chi.Router) {
			// Get self
			r.Get("/", getOrg)
			// create a repository
			r.Post("/repo", createRepo)
			// create a project
			r.Post("/project", createProject)
			// create a folder
			r.Post("/folder", createFolder)
		})
	})

	router.Route("/folder", func(r chi.Router) {
		r.Route("/{folderId}", func(r chi.Router) {
			r.Get("/", getFolder)
			// create a folder
			r.Post("/folder", createFolder)
			// create a project
			// 	requirement: No project should exist within the path
			// 	eg. /workspace/project/folder/folder/project is invalid
			//  but /workspace/folder/folder/project is valid.
			r.Post("/project", createProject)
		})
	})

	router.Route("/project", func(r chi.Router) {
		r.Route("/{projectId}", func(r chi.Router) {
			r.Get("/", getProject)
		})
	})

	router.Route("/repo", func(r chi.Router) {
		r.Get("/{repoId}", getRepo)
	})
}
