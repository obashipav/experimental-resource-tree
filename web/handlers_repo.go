package web

import (
	"encoding/json"
	"github.com/OBASHITechnology/resourceList/DB"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"net/http"
)

/*
	Repository Handlers
*/
func createRepo(w http.ResponseWriter, r *http.Request) {
	path := models.CleanSlashFromPath(models.ExtractParentPath(r.URL.Path, repo.APIRoute))

	decoder := json.NewDecoder(r.Body)
	content := repo.CreateRequest{BaseInfo: models.BaseInfo{PreviousURL: path}}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a repo resource", http.StatusBadRequest)
		return
	}

	// Validation
	content.Valid()

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
	//path := chi.URLParam(r, repo.URLParam)
	path := models.CleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetRepo(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the repo resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}
