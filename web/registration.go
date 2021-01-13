package web

import (
	"fmt"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"github.com/OBASHITechnology/resourceList/models/org"
	"github.com/OBASHITechnology/resourceList/models/project"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
)

var (
	Engine *chi.Mux

	AllowAPI jsoniter.API

	relaxedAPI = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  false,
	}.Froze()
)

func init() {
	AllowAPI = relaxedAPI

	Engine = chi.NewRouter()
	Engine.Use(middleware.Logger)
	Engine.Use(middleware.Recoverer)
}

func Registration() *chi.Mux {
	Engine.Route(org.APIRoute, func(r chi.Router) {
		r.Post("/", createOrg)

		r.Route(fmt.Sprintf("/{%s}", org.URLParam), func(r chi.Router) {
			// Get self
			r.Get("/", getOrg)
			// create a repository
			r.Post(repo.APIRoute, createRepo)
			// create a project
			r.Post(project.APIRoute, createProject)
			// create a folder
			r.Post(folder.APIRoute, createFolder)
		})
	})

	Engine.Route(folder.APIRoute, func(r chi.Router) {
		r.Route(fmt.Sprintf("/{%s}", folder.URLParam), func(r chi.Router) {
			r.Get("/", getFolder)
			// create a folder
			r.Post(folder.APIRoute, createFolder)
			// create a project
			// 	requirement: No project should exist within the path
			// 	eg. /workspace/project/folder/folder/project is invalid
			//  but /workspace/folder/folder/project is valid.
			r.Post(project.APIRoute, createProject)

			// Deleting a folder can be tricky, risky and dangerous
			//  however we should inform the user that the action
			//  is final.
			r.Delete("/", deleteFolder)
		})
	})

	Engine.Route(project.APIRoute, func(r chi.Router) {
		r.Route(fmt.Sprintf("/{%s}", project.URLParam), func(r chi.Router) {
			r.Get("/", getProject)
			// create a folder
			r.Post(folder.APIRoute, createFolder)
		})
	})

	Engine.Route(repo.APIRoute, func(r chi.Router) {
		r.Get(fmt.Sprintf("/{%s}", repo.URLParam), getRepo)
	})
	return Engine
}
