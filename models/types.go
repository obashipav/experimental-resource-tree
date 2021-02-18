package models

import (
	"errors"
	"strings"
)

type (
	Timestamp struct {
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}

	UserAction struct {
		Owner     string `json:"owner"`
		UpdatedBy string `json:"updatedBy"`
	}

	BaseInfo struct {
		Label       string `json:"label"`
		Description string `json:"description"`
		AltLabel    string `json:"altLabel"`
		PreviousURL string `json:"previousURL,omitempty"`
	}

	ResourceInfo struct {
		BaseInfo
		UserAction
		Timestamp
	}

	NextURLs map[string]*ResourceInfo

	CreateResponse struct {
		ResourceID  string   `json:"-"`
		URL         string   `json:"selfUrl"`
		PreviousURL string   `json:"previousURL,omitempty"`
		NextURLs    NextURLs `json:"nextURLs"`
	}

	ResourceHistory struct {
		Timestamp
		UserAction
	}

	/*HierarchyInfo struct {
		Type  string `json:"type"`  // this should point to the actual table
		Order int    `json:"order"` // the hierarchy order of the map
	}

	// the key is the plain resource id as a uuid
	//  using the type and the key we can retrieve the resource from the table.
	HierarchyMap map[string]*HierarchyInfo*/

	Hierarchy struct {
		List []string `json:"list"`
	}
)

// AssignOwnerWhenCreating this function assigns the owner value to the updated by field only when we
//  create resources. This function should be only used when creating things.
func (u *UserAction) AssignOwnerWhenCreating() {
	u.UpdatedBy = u.Owner
}

func (b *BaseInfo) CleanLabels() {
	b.Label = TrimSpacesInBetween(strings.TrimSpace(b.Label))
	b.AltLabel = TrimSpacesInBetween(strings.TrimSpace(b.AltLabel))
	b.Description = TrimSpacesInBetween(strings.TrimSpace(b.Description))
}

// AddResource function accepts two parameters which must be the previous and current path aliases.
//  The function performs some check before adding the new alias to the path. Simply returns an error.
func (h *Hierarchy) AddResource(previous, aliasPath string) error {
	if h == nil {
		return errors.New("the hierarchy is not set, please initiate the hierarchy")
	}

	// the aliasPath must be of type {type}/{aliasId}
	resourceType := strings.Split(aliasPath, "/")[0]
	// we will store the values of the slice to a map for quick access during the checks
	localMap := make(map[string]int)
	for i, v := range h.List {
		localMap[v] = i + 1
	}

	/*
		Error checks
	*/
	switch {
	case len(h.List)==0 && resourceType != "org":
		return errors.New("the root resource must be an organisation")
	case len(h.List) > 0 && resourceType == "org":
		return errors.New("the organisation must be only at the root")
	case len(h.List) > 0 && resourceType != "org" && h.List[len(h.List)-1] != previous :
		return errors.New("the previous resource doesn't match with the path")
	case len(h.List) > 0 && resourceType != "org" && h.List[len(h.List)-1] == aliasPath :
		return errors.New("the resource is already part of the path")
	default:
		h.List = append(h.List, aliasPath)
	}
	return nil
}
