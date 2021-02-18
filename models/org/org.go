package org

import (
	"github.com/OBASHITechnology/resourceList/models"
)

const (
	APIRoute  = "/org"
	DBTable   = "org"
	URIScheme = "org"
	URLParam  = "orgId"
)

type (
	CreateRequest struct {
		ID    string `json:"-"`
		Alias string `json:"-"`
		models.BaseInfo
		models.UserAction
		// mainly used internally.
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
	if c.Hierarchy.List == nil {
		c.Hierarchy.List = make([]string,0)
	}
}
