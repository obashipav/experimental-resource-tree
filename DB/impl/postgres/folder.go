package postgres

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/folderdb"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/pathdb"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"github.com/OBASHITechnology/resourceList/models/path"
	"github.com/OBASHITechnology/resourceList/util/shortID"
	"log"
)

func (s *store) CreateFolder(request *folder.CreateRequest) (*models.CreateResponse, error) {
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

	parent, err = pathdb.GetPathDetails(tx, request.PreviousURL)
	if err != nil {
		return nil, err
	}

	//request.ID = uuid.NewID()
	request.PathURI = shortID.NewWithURL(request.PreviousURL)
	request.HierarchyMap = parent.Hierarchy
	err = request.AddResource(parent.ResourceID, request.ID, folder.DBTable)
	if err != nil {
		return nil, err
	}

	var response *models.CreateResponse
	response, err = folderdb.Create(tx, request)
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

func (s *store) GetFolder(url string) (*folder.GetResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	tx, err := s.pool.Begin(ctx)
	defer common.HandleTransactionRollback(tx, ctx, cancel)
	if err != nil {
		log.Println("there was an error with the tx begin: ", err)
		return nil, err
	}

	var response *folder.GetResponse
	response, err = folderdb.Get(tx, url)
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

// DeleteFolder will delete the resource if there are no dependencies. Otherwise, if the kind of deletion is soft and
//  and there are dependencies, an error message will be returned to the user.
func (s *store) DeleteFolder(url string, force bool) error {
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	tx, err := s.pool.Begin(ctx)
	defer common.HandleTransactionRollback(tx, ctx, cancel)
	if err != nil {
		log.Println("there was an error with the tx begin: ", err)
		return err
	}

	var resource = &path.GetResponse{}
	resource, err = pathdb.GetPathDetails(tx, url)
	if err != nil {
		return err
	}
	resource.NextURLs, err = pathdb.GetNextURLs(tx, url)
	if err != nil {
		return err
	}

	// delete the Resource

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("there was an error with the tx commit: ", err)
		return err
	}

	return nil
}