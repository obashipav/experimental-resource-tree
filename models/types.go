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

	HierarchyInfo struct {
		Type  string `json:"type"`  // this should point to the actual table
		Order int    `json:"order"` // the hierarchy order of the map
	}

	// the key is the plain resource id as a uuid
	//  using the type and the key we can retrieve the resource from the table.
	HierarchyMap map[string]*HierarchyInfo
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

// GetHierarchyAsPath returns the uuid hierarchy as a uri style path, separated by front slash
//  {root uuid} / { resource uuid} / ... / { resource uuid}
//  The order is important to return the right path.
//
//  This function is not safe for now, because it considers the order to be correct.
func (h HierarchyMap) GetHierarchyAsPath() string {
	if len(h) == 0 {
		return ""
	}
	list := make([]string, len(h))
	for k, v := range h {
		if v.Order >= len(list) {
			return ""
		}
		list[v.Order] = k
	}
	return strings.Join(list, "/")
}

func (h HierarchyMap) AddResource(previous, id, resourceType string) error {
	if h == nil {
		return errors.New("the hierarchy is not set, check again")
	}
	if _, exists := h[previous]; !exists && resourceType != "org" {
		return errors.New("the parent doesn't exist in the path, check again")
	}

	if _, exists := h[id]; exists {
		return errors.New("the resource is already part of the hierarchy, check again")
	}

	if resourceType == "org" {
		h[id] = &HierarchyInfo{
			Type:  resourceType,
			Order: 0,
		}
	} else {
		h[id] = &HierarchyInfo{
			Type:  resourceType,
			Order: h[previous].Order + 1,
		}
	}
	return nil
}
