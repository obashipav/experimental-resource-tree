package postgres

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/pathdb"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/repodb"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/path"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"log"
)

func (s *store) CreateRepo(request *repo.CreateRequest) (*models.CreateResponse, error) {
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
	response, err = repodb.Create(tx, request)
	if err != nil {
		return nil, err
	}

	parent.Hierarchy[response.ResourceID] = &path.ResourceInfo{
		Type:  "repo",
		Order: parent.Hierarchy[parent.ResourceID].Order + 1,
	}

	response.URL, err = pathdb.AddResource(tx, &path.CreateRequest{
		ResourceID:  response.ResourceID,
		PreviousURL: request.PathURL,
		Type:        "repo",
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

func (s *store) GetRepo(url string) (*repo.GetResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	tx, err := s.pool.Begin(ctx)
	defer common.HandleTransactionRollback(tx, ctx, cancel)
	if err != nil {
		log.Println("there was an error with the tx begin: ", err)
		return nil, err
	}

	var response *repo.GetResponse
	response, err = repodb.Get(tx, url)
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
