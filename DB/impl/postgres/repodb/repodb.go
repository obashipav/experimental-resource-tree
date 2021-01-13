package repodb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"log"
)

func Create(db common.QueryRower, request *repo.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.repository (path_uri, previous_uri, type, hierarchy, label, alt_label, description, owner, updated_by) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	err := db.QueryRow(ctx, query, request.Alias, request.PreviousURL, repo.DBTable, request.HierarchyMap, request.Label, request.AltLabel, request.Description, request.Owner, request.UpdatedBy).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the repo table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response, URL: models.GetRealPath(request.Alias), PreviousURL: models.GetRealPath(request.PreviousURL)}, nil
}

func Get(db common.QueryRower, url string) (*repo.GetResponse, error) {
	const query = `select previous_uri, label, alt_label, description, owner, updated_by, created_at, updated_at
		from demo.base
		where path_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &repo.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealPath(url)}}
	err := db.QueryRow(ctx, query, url).Scan(&response.Path.PreviousURL, &response.Label, &response.AltLabel, &response.Description, &response.History.Owner,
		&response.History.UpdatedBy, &response.History.CreatedAt, &response.History.UpdatedAt)
	if err != nil {
		log.Println("failed to select the repo", err)
		return nil, err
	}
	response.Path.PreviousURL = models.GetRealPath(response.Path.PreviousURL)
	return response, err
}
