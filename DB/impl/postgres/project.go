package postgres

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/pathdb"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/projectdb"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/path"
	"github.com/OBASHITechnology/resourceList/models/project"
	"log"
)

func (s *store) CreateProject(request *project.CreateRequest) (*models.CreateResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	tx, err := s.pool.Begin(ctx)
	defer common.HandleTransactionRollback(tx, ctx, cancel)
	if err != nil {
		log.Println("there was an error with the tx begin: ", err)
		return nil, err
	}

	// first check if the path is valid
	//  leave the permissions for later
	var parent *path.GetResponse
	parent, err = pathdb.GetPathDetails(tx, request.PathURL)
	if err != nil {
		return nil, err
	}

	var response *models.CreateResponse
	response, err = projectdb.Create(tx, request)
	if err != nil {
		return nil, err
	}

	parent.Hierarchy[response.ResourceID] = &path.ResourceInfo{
		Type:  "project",
		Order: parent.Hierarchy[parent.ResourceID].Order + 1,
	}

	response.URL, err = pathdb.AddResource(tx, &path.CreateRequest{
		ResourceID:  response.ResourceID,
		PreviousURL: request.PathURL,
		Type:        "project",
		Hierarchy:   parent.Hierarchy,
	})

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("there was an error with the tx commit: ", err)
		return nil, err
	}

	return response, nil
}

func (s *store) GetProject(url string) (*project.GetResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	tx, err := s.pool.Begin(ctx)
	defer common.HandleTransactionRollback(tx, ctx, cancel)
	if err != nil {
		log.Println("there was an error with the tx begin: ", err)
		return nil, err
	}

	var response *project.GetResponse
	response, err = projectdb.Get(tx, url)
	if err != nil {
		return nil, err
	}

	response.Path.NextURLs, err = pathdb.GetNextURLs(tx, url)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("there was an error with the tx commit: ", err)
		return nil, err
	}

	return response, nil
}
