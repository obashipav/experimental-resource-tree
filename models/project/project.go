package project

import "github.com/OBASHITechnology/resourceList/models"

const (
	APIRoute  = "/project"
	DBTable   = "project"
	URIScheme = "project"
	URLParam  = "projectId"
)

type (
	CreateRequest struct {
		ID    string `json:"-"`
		Alias string `json:"-"`
		Color string `json:"color"`
		models.BaseInfo
		models.UserAction
		models.HierarchyMap
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
