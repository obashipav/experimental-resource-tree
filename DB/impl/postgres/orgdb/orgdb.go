package orgdb

import (
	"context"
	"fmt"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/org"
	"log"
)

func Create(db common.QueryRower, request *org.CreateRequest) (*models.CreateResponse, error) {
	const query = `insert into demo.org (path_uri, previous_uri, type, hierarchy, label, alt_label, description, owner, updated_by) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response string
	err := db.QueryRow(ctx, query, request.PathURI, request.PreviousURL, org.DBTable, request.HierarchyMap, request.Label, request.AltLabel, request.Description, request.Owner, request.UpdatedBy).Scan(&response)
	if err != nil {
		log.Println("failed to insert into the org table: ", err)
		return nil, err
	}
	return &models.CreateResponse{ResourceID: response, URL: models.GetRealURL(fmt.Sprintf("%s%s", org.URIScheme, request.PathURI))}, nil
}

func Get(db common.QueryRower, url string) (*org.GetResponse, error) {
	const query = `select label, alt_label, description, owner, updated_by, created_at, updated_at
		from demo.base
		where path_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = &org.GetResponse{History: models.ResourceHistory{}, Path: models.CreateResponse{URL: models.GetRealURL(fmt.Sprintf("%s%s", org.URIScheme, url))}}
	err := db.QueryRow(ctx, query, url).Scan(&response.Label, &response.AltLabel, &response.Description, &response.History.Owner, &response.History.UpdatedBy, &response.History.CreatedAt, &response.History.UpdatedAt)
	if err != nil {
		log.Println("failed to select the org", err)
		return nil, err
	}
	return response, err
}
