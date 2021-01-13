package pathdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/path"
	"github.com/OBASHITechnology/resourceList/util/shortID"
	"log"
)

// Deprecated:
func AddResource(db common.QueryRower, resource *path.CreateRequest) (string, error) {
	const query = `insert into demo.respath (path_uri, id, type, previous_uri, hierarchy) values ($1, $2, $3, $4, $5) on conflict (path_uri) do nothing returning path_uri;`
	var ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
	defer cancel()

	// not ready, I need to prepend the URI of the type
	pathURL := shortID.NewWithURL(fmt.Sprintf(resource.Hierarchy.GetHierarchyAsPath(), resource.ResourceID))

	var response string
	err := db.QueryRow(ctx, query, pathURL, resource.ResourceID, resource.Type, resource.PreviousURL, resource.Hierarchy).Scan(&response)
	if err != nil {
		log.Println("failed to insert the resource to the path: ", err)
		return "", err
	}

	return models.GetRealURL(response), nil
}

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
		response[models.GetRealURL(nextUrl)] = &next
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
	if _, exists := response.Hierarchy[response.ResourceID]; !exists || len(response.Hierarchy) == 0 {
		log.Println("failed to resolve the hierarchy path", response)
		return nil, errors.New("failed to return the hierarchy")
	}
	return response, nil
}
