package repo

import (
	"github.com/OBASHITechnology/resourceList/models"
)

const (
	APIRoute  = "/repository"
	DBTable   = "repository"
	URIScheme = "repository"
	URLParam  = "repoId"
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
