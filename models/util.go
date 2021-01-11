package models

import "fmt"

func GetRealURL(path string) string {
	return fmt.Sprintf("http://localhost:8080/%s", path)
}
