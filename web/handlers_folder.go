package web

import (
	"encoding/json"
	"github.com/OBASHITechnology/resourceList/DB"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"net/http"
)

/*
	Folder
*/
func createFolder(w http.ResponseWriter, r *http.Request) {
	path := models.CleanSlashFromPath(models.ExtractParentPath(r.URL.Path, folder.APIRoute))

	decoder := json.NewDecoder(r.Body)
	content := folder.CreateRequest{BaseInfo: models.BaseInfo{PreviousURL: path}}
	err := decoder.Decode(&content)
	if err != nil {
		httpAbortWithMessage(w, "failed to create a folder resource", http.StatusBadRequest)
		return
	}

	// Validation
	content.Valid()

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
	//path := chi.URLParam(r, folder.URLParam)
	path := models.CleanSlashFromPath(r.URL.Path)

	response, err := DB.Store.GetFolder(path)
	if err != nil {
		httpAbortWithMessage(w, "failed to get the folder resource", http.StatusInternalServerError)
		return
	}
	httpResponse(w, response, http.StatusOK)
}

func deleteFolder(w http.ResponseWriter, r *http.Request) {
	//path := models.CleanSlashFromPath(r.URL.Path)
	//kind := r.URL.Query().Get(folder.DeleteQuery)
}
