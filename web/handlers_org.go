package web

import (
	"encoding/json"
	"github.com/OBASHITechnology/resourceList/DB"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/org"
	"net/http"
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

	// Validation
	content.Valid()

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
	// path := chi.URLParam(r, org.URLParam)
	path := models.CleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetOrg(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the org resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}
