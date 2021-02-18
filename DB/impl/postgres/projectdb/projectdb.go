package projectdb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/project"
	"log"
)

func Create(db common.QueryRower, request *project.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.project (path_uri, previous_uri, type, hierarchy, hlevel, label, alt_label, description, color, owner, updated_by) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	err := db.QueryRow(ctx, query, request.Alias, request.PreviousURL, project.DBTable, request.Hierarchy.List, len(request.Hierarchy.List), request.Label, request.AltLabel, request.Description, request.Color, request.Owner, request.UpdatedBy).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the project table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response, URL: models.GetRealPath(request.Alias), PreviousURL: models.GetRealPath(request.PreviousURL)}, nil
}

func Get(db common.QueryRower, url string) (*project.GetResponse, error) {
	const query = `select previous_uri, label, alt_label, description, owner, updated_by, created_at, updated_at
		from demo.base
		where path_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &project.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealPath(url)}}
	err := db.QueryRow(ctx, query, url).Scan(&response.Path.PreviousURL, &response.Label, &response.AltLabel, &response.Description, &response.History.Owner,
		&response.History.UpdatedBy, &response.History.CreatedAt, &response.History.UpdatedAt)
	if err != nil {
		log.Println("failed to select the project", err)
		return nil, err
	}
	response.Path.PreviousURL = models.GetRealPath(response.Path.PreviousURL)
	return response, err
}
