package projectdb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/project"
	"log"
	"time"
)

func Create(db common.QueryRower, request *project.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.project (label, path, created_by, updated_by, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	var timestamp = time.Now().UnixNano() / 1000000
	err := db.QueryRow(ctx, query, request.Label, request.PathURL, request.CreatedBy, request.UpdatedBy, timestamp, timestamp).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the repo table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response}, nil
}

func Get(db common.QueryRower, url string) (*project.GetResponse, error) {
	const query = `select p.label, p.created_by, p.created_at, p.updated_by, p.updated_at, coalesce(r.previous_url,'')
		from only demo.project p
		inner join demo.respath r on r.resource_id = p.id
		where r.path_url = $1 and not p.deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &project.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealURL(url)}}
	err := db.QueryRow(ctx, query, url).Scan(
		&response.Label, &response.History.CreatedBy, &response.History.CreatedAt,
		&response.History.UpdatedBy, &response.History.UpdatedAt, &response.Path.PreviousURL)
	if err != nil {
		log.Println("failed to select the repo", err)
		return nil, err
	}
	response.Path.PreviousURL = models.GetRealURL(response.Path.PreviousURL)
	return response, err
}
