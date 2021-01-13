package path

import (
	"github.com/OBASHITechnology/resourceList/models"
)

const (
	DBTable = "base"
)

type (
	CreateRequest struct {
		ResourceID  string              `json:"-"`           // the resource id as a uuid.
		Type        string              `json:"type"`        // points to the resource table
		PreviousURL string              `json:"previousURL"` // parent path as a URL link
		Hierarchy   models.HierarchyMap `json:"hierarchy"`
	}

	GetResponse struct {
		ResourceID  string              `json:"-"`
		URL         string              `json:"selfURL"`
		Type        string              `json:"type"`
		PreviousURL string              `json:"previousURL"`
		Hierarchy   models.HierarchyMap `json:"hierarchy"`
		models.NextURLs
	}
)
