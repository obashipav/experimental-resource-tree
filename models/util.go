package models

import (
	"fmt"
	"strings"
)

const serverSchema = `http://localhost:8080/%s`

// GetRealURL function prepends the server's http address, in order to form a proper URL.
func GetRealURL(path string) string {
	return fmt.Sprintf(serverSchema, path)
}

// CleanSlashFromPath function removes any trailing slash from the relative path. This will allow us to
//  start doing queries to the db without having problems, (the problems will be `Not Found`).
func CleanSlashFromPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return strings.TrimPrefix(path, "/")
	}
	return path
}

// ExtractParentPath function is mainly used for cleaning the URL Path from suffixes, in order to
//  extract the path and starting querying the db.
//   EG.  http://localhost:8080/org/{some random string}/project
//   looking at the URL above, we only care for the `org/{some random string}`, the chi engine returns the
//   the relative path with the suffix, which we want to trim.
func ExtractParentPath(path, suffix string) string {
	if strings.HasSuffix(path, suffix) {
		return strings.TrimSuffix(path, suffix)
	}
	return path
}

// TrimSpacesInBetween this function is designed to clean tabs and spaces from labels and descriptions.
func TrimSpacesInBetween(label string) string {
	if label == "" {
		return label
	}
	source := strings.Split(label, " ")
	if len(source) == 0 {
		return label
	}

	var target = make([]string, 0)
	for _, word := range source {
		if word != "" {
			target = append(target, word)
		}
	}
	return strings.Join(target, " ")
}
