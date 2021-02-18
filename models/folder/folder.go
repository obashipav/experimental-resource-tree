package folder

import "github.com/OBASHITechnology/resourceList/models"

const (
	// API route
	APIRoute = "/folder"
	// URI query keys
	DeleteQuery = "kind"
	// Resource DB table
	DBTable   = "folder"
	URIScheme = "folder"
	URLParam  = "folderId"
)

type (
	CreateRequest struct {
		ID    string `json:"-"`
		Alias string `json:"-"`
		models.BaseInfo
		models.UserAction
		// mainly used internally
		Hierarchy models.Hierarchy `json:"-"`
	}

	GetResponse struct {
		History models.ResourceHistory `json:"history"`
		Path    models.CreateResponse  `json:"path"`
		models.BaseInfo
	}
)

func (c *CreateRequest) Valid() {
	c.BaseInfo.CleanLabels()
	c.UserAction.AssignOwnerWhenCreating()
}
