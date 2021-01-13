package web

import (
	"encoding/json"
	"github.com/OBASHITechnology/resourceList/DB"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/project"
	"net/http"
)

/*
	Project
*/
func createProject(w http.ResponseWriter, r *http.Request) {
	path := models.CleanSlashFromPath(models.ExtractParentPath(r.URL.Path, project.APIRoute))

	decoder := json.NewDecoder(r.Body)
	content := project.CreateRequest{BaseInfo: models.BaseInfo{PreviousURL: path}}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a project resource", http.StatusBadRequest)
		return
	}

	// Validation
	content.Valid()

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
	//path := chi.URLParam(r, project.URLParam)
	path := models.CleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetProject(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the project resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}
