package org

import (
	"github.com/OBASHITechnology/resourceList/models"
)

type (
	CreateRequest struct {
		Label string `json:"label"`
		models.UserAction
	}

	GetResponse struct {
		Label   string                 `json:"label"`
		History models.ResourceHistory `json:"history"`
		Path    models.CreateResponse  `json:"path"`
	}
)
