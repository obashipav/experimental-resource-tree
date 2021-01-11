package orgdb

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/org"
	"log"
	"time"
)

func Create(db common.QueryRower, request *org.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.org (label, created_by, updated_by, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	var timestamp = time.Now().UnixNano() / 1000000
	err := db.QueryRow(ctx, query, request.Label, request.CreatedBy, request.UpdatedBy, timestamp, timestamp).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the org table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response}, nil
}

func Get(db common.QueryRower, url string) (*org.GetResponse, error) {
	const query = `select o.label, o.created_by, o.created_at, o.updated_by, o.updated_at, coalesce(p.previous_url,'')
		from only demo.org o
		inner join demo.respath p on p.resource_id = o.id
		where p.path_url = $1 and not o.deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &org.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealURL(url)}}
	err := db.QueryRow(ctx, query, url).Scan(&response.Label, &response.History.CreatedBy, &response.History.CreatedAt,
		&response.History.UpdatedBy, &response.History.UpdatedAt, &response.Path.PreviousURL)
	if err != nil {
		log.Println("failed to select the org", err)
		return nil, err
	}
	return response, err
}
