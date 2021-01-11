package models

type (
	CreateResponse struct {
		ResourceID  string   `json:"-"`
		URL         string   `json:"selfUrl"`
		PreviousURL string   `json:"previousURL,omitempty"`
		NextURLs    []string `json:"nextURLs"`
	}

	Timestamp struct {
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}

	UserAction struct {
		CreatedBy string `json:"createdBy"`
		UpdatedBy string `json:"updatedBy"`
	}

	ResourceHistory struct {
		Timestamp
		UserAction
	}
)
