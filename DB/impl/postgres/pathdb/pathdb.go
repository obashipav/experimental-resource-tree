package pathdb

import (
	"context"
	"errors"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/path"
	"log"
)

func GetNextURLs(db common.Queryer, url string) (models.NextURLs, error) {
	const query = `select path_uri, label, alt_label, description, owner, updated_by, created_at, updated_at from demo.base where previous_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()
	var response = make(models.NextURLs)
	rows, err := db.Query(ctx, query, url)
	if err != nil {
		log.Println("failed to find the next urls: ", err)
		return nil, err
	}
	for rows.Next() {
		var next = models.ResourceInfo{}
		var nextUrl string
		err = rows.Scan(&nextUrl, &next.Label, &next.AltLabel, &next.Description, &next.Owner, &next.UpdatedBy, &next.CreatedAt, &next.UpdatedAt)
		if err != nil {
			log.Println("failed to scan the url: ", err)
			return nil, err
		}
		response[models.GetRealPath(nextUrl)] = &next
	}
	return response, nil
}

func GetPathDetails(db common.QueryRower, url string) (*path.GetResponse, error) {
	const query = `select path_uri, id, type, coalesce(previous_uri, ''), hierarchy 
		from demo.base where path_uri = $1 and not deleted;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	var response = &path.GetResponse{}
	err := db.QueryRow(ctx, query, url).Scan(&response.URL, &response.ResourceID, &response.Type,
		&response.PreviousURL, &response.Hierarchy)
	if err != nil {
		log.Println("failed to scan the base table: ", err)
		return nil, err
	}
	// quality check
	if _, exists := response.Hierarchy[response.URL]; !exists || len(response.Hierarchy) == 0 {
		log.Println("failed to resolve the hierarchy path", response)
		return nil, errors.New("failed to return the hierarchy")
	}
	return response, nil
}
