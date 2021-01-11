package repodb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"log"
	"time"
)

func Create(db common.QueryRower, request *repo.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.repo (label, path, created_by, updated_by, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) returning id;`
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

func Get(db common.QueryRower, url string) (*repo.GetResponse, error) {
	const query = `select r.label, r.created_by, r.created_at, r.updated_by, r.updated_at, coalesce(p.previous_url,'')
		from only demo.repo r
		inner join demo.respath p on p.resource_id = r.id
		where p.path_url = $1 and not r.deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &repo.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealURL(url)}}
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
