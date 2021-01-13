package core

import (
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"github.com/OBASHITechnology/resourceList/models/org"
	"github.com/OBASHITechnology/resourceList/models/project"
	"github.com/OBASHITechnology/resourceList/models/repo"
)

type (
	IStore interface {
		CreateOrg(request *org.CreateRequest) (*models.CreateResponse, error)
		GetOrg(url string) (*org.GetResponse, error)

		CreateProject(request *project.CreateRequest) (*models.CreateResponse, error)
		GetProject(url string) (*project.GetResponse, error)
		//
		CreateFolder(request *folder.CreateRequest) (*models.CreateResponse, error)
		GetFolder(url string) (*folder.GetResponse, error)
		DeleteFolder(url string, force bool) error

		CreateRepo(request *repo.CreateRequest) (*models.CreateResponse, error)
		GetRepo(url string) (*repo.GetResponse, error)
	}
)
