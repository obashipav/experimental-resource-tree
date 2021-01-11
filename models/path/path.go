package path

import "strings"

type (
	ResourceInfo struct {
		Type  string `json:"type"`  // this should point to the actual table
		Order int    `json:"order"` // the hierarchy order of the map
	}

	// the key is the plain resource id as a uuid
	//  using the type and the key we can retrieve the resource from the table.
	HierarchyMap map[string]*ResourceInfo

	CreateRequest struct {
		ResourceID  string       `json:"resourceId"`  // the resource id as a uuid.
		Type        string       `json:"type"`        // points to the resource table
		PreviousURL string       `json:"previousURL"` // parent path as a URL link
		Hierarchy   HierarchyMap `json:"hierarchy"`
	}

	GetResponse struct {
		ResourceID  string       `json:"-"`
		URL         string       `json:"selfURL"`
		Type        string       `json:"type"`
		PreviousURL string       `json:"previousURL"`
		Hierarchy   HierarchyMap `json:"hierarchy"`
		NextURLs    []string     `json:"nextURLs"`
	}
)

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
