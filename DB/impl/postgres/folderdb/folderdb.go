package folderdb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"log"
	"time"
)

func Create(db common.QueryRower, request *folder.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.folder (label, path, created_by, updated_by, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) returning id;`
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

func Get(db common.QueryRower, url string) (*folder.GetResponse, error) {
	const query = `select f.label, f.created_by, f.created_at, f.updated_by, f.updated_at, coalesce(r.previous_url,'')
		from only demo.folder f
		inner join demo.respath r on r.resource_id = f.id
		where r.path_url = $1 and not f.deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &folder.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealURL(url)}}
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
