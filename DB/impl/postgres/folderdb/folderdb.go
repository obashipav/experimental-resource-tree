package folderdb

import (
	"context"
	"errors"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"log"
)

func Create(db common.QueryRower, request *folder.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.folder (path_uri, previous_uri, type, hierarchy, hlevel, label, alt_label, description, owner, updated_by) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	err := db.QueryRow(ctx, query, request.Alias, request.PreviousURL, folder.DBTable, request.Hierarchy.List, len(request.Hierarchy.List), request.Label, request.AltLabel, request.Description, request.Owner, request.UpdatedBy).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the folder table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response, URL: models.GetRealPath(request.Alias), PreviousURL: models.GetRealPath(request.PreviousURL)}, nil
}

func Get(db common.QueryRower, url string) (*folder.GetResponse, error) {
	const query = `select previous_uri, label, alt_label, description, owner, updated_by, created_at, updated_at
		from demo.base
		where path_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &folder.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealPath(url)}}
	err := db.QueryRow(ctx, query, url).Scan(&response.Path.PreviousURL, &response.Label, &response.AltLabel, &response.Description, &response.History.Owner,
		&response.History.UpdatedBy, &response.History.CreatedAt, &response.History.UpdatedAt)
	if err != nil {
		log.Println("failed to select the folder", err)
		return nil, err
	}
	response.Path.PreviousURL = models.GetRealPath(response.Path.PreviousURL)
	return response, err
}

func Delete(db common.Execer, url string) error {
	const query = `update demo.folder set deleted = true where uri = $1;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	tag, err := db.Exec(ctx, query, url)
	if err != nil {
		log.Println("failed to delete the folder: ", err)
		return err
	}
	if tag.RowsAffected() == 0 {
		log.Println("nothing has been deleted, url: ", url)
		return errors.New("not found")
	}
	return nil
}
