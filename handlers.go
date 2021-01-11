package main

import (
	"encoding/json"
	"github.com/OBASHITechnology/resourceList/DB"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"github.com/OBASHITechnology/resourceList/models/org"
	"github.com/OBASHITechnology/resourceList/models/project"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"net/http"
	"strings"
)

/*
	Organisation Handlers
*/
func createOrg(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	content := org.CreateRequest{}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create an org resource", http.StatusBadRequest)
		return
	}

	// postgres
	var response *models.CreateResponse
	response, err = DB.Store.CreateOrg(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create an org resource", http.StatusInternalServerError)
		return
	}

	// return
	httpResponse(w, response, http.StatusCreated)
}

func getOrg(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetOrg(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the org resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}

func cleanSlashFromPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return strings.TrimPrefix(path, "/")
	}
	return path
}

/*
	Repository Handlers
*/
func createRepo(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(extractParentPath(r.URL.Path, "/repo"))

	decoder := json.NewDecoder(r.Body)
	content := repo.CreateRequest{PathURL: path}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a repo resource", http.StatusBadRequest)
		return
	}

	// postgres
	var response *models.CreateResponse
	response, err = DB.Store.CreateRepo(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a repo resource", http.StatusInternalServerError)
		return
	}

	// return
	httpResponse(w, response, http.StatusCreated)
}

func getRepo(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetRepo(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the repo resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}

func extractParentPath(path, suffix string) string {
	if strings.HasSuffix(path, suffix) {
		return strings.TrimSuffix(path, suffix)
	}
	return path
}

/*
	Project
*/
func createProject(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(extractParentPath(r.URL.Path, "/project"))

	decoder := json.NewDecoder(r.Body)
	content := project.CreateRequest{PathURL: path}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a project resource", http.StatusBadRequest)
		return
	}

	// postgres
	var response *models.CreateResponse
	response, err = DB.Store.CreateProject(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a project resource", http.StatusInternalServerError)
		return
	}

	// return
	httpResponse(w, response, http.StatusCreated)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetProject(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the project resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}

/*
	Folder
*/
func createFolder(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(extractParentPath(r.URL.Path, "/folder"))

	decoder := json.NewDecoder(r.Body)
	content := folder.CreateRequest{PathURL: path}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a folder resource", http.StatusBadRequest)
		return
	}

	// postgres
	var response *models.CreateResponse
	response, err = DB.Store.CreateFolder(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a folder resource", http.StatusInternalServerError)
		return
	}

	// return
	httpResponse(w, response, http.StatusCreated)
}

func getFolder(w http.ResponseWriter, r *http.Request) {
	path := cleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetFolder(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the folder resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}
