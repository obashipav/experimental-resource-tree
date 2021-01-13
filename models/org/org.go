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
	if c.HierarchyMap == nil {
		c.HierarchyMap = make(models.HierarchyMap)
	}
}
