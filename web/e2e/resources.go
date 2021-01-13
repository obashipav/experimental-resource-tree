package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/OBASHITechnology/resourceList/models"
	"github.com/OBASHITechnology/resourceList/models/folder"
	"github.com/OBASHITechnology/resourceList/models/org"
	"github.com/OBASHITechnology/resourceList/models/project"
	"github.com/OBASHITechnology/resourceList/models/repo"
	"github.com/OBASHITechnology/resourceList/web"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

type Pyramid struct {
	Owner        string                 `json:"owner"`
	Org          *org.GetResponse       `json:"root"`
	Folders      []*folder.GetResponse  `json:"folders"`
	Projects     []*project.GetResponse `json:"projects"`
	Repositories []*repo.GetResponse    `json:"repositories"`
}

var createResources = func() {

	const (
		// owners
		PavlosUser = "Pavlos"

		// address
		URIScheme = "http://localhost:8080/"
	)

	var (
		payload  = []byte("")
		err      error
		content  models.CreateResponse
		recorder *httptest.ResponseRecorder

		// pyramids
		pavlosPyramid = Pyramid{
			Owner:        PavlosUser,
			Org:          nil,
			Folders:      make([]*folder.GetResponse, 0),
			Projects:     make([]*project.GetResponse, 0),
			Repositories: make([]*repo.GetResponse, 0),
		}
	)

	ginkgo.Context("As a normal user", func() {

		ginkgo.Context("who wants to create a root resource named org", func() {

			ginkgo.When("I make a valid create request", func() {
				ginkgo.It("should create the org and return the link", func() {
					payload, err = json.Marshal(&org.CreateRequest{
						BaseInfo:   models.BaseInfo{Label: "Forth Valley College", Description: "Our first college"},
						UserAction: models.UserAction{Owner: PavlosUser},
					})
					gomega.Expect(err).To(gomega.BeNil())

					recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%sorg", URIScheme), payload).GetRecorder(web.Engine)
					gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

					err = json.Unmarshal(recorder.Body.Bytes(), &content)
					gomega.Expect(err).To(gomega.BeNil())

					recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
					gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

					var response org.GetResponse
					err = json.Unmarshal(recorder.Body.Bytes(), &response)
					gomega.Expect(err).To(gomega.BeNil())

					pavlosPyramid.Org = &response
				})
			})
		})

		ginkgo.Context("using this org", func() {

			ginkgo.Context("who wants to create repositories", func() {
				ginkgo.When("I make a valid create request", func() {
					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&repo.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Dataflow Unit Assets", Description: "Public"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/repository", pavlosPyramid.Org.Path.URL), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response repo.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Repositories = append(pavlosPyramid.Repositories, &response)
					})
				})

				ginkgo.When("I make a second create request for repo", func() {
					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&repo.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Falkirk Assets", Description: "Needs access"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/repository", pavlosPyramid.Org.Path.URL), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response repo.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Repositories = append(pavlosPyramid.Repositories, &response)
					})
				})
			})

			ginkgo.Context("who wants to create projects", func() {
				ginkgo.When("I make a valid create request", func() {

					var tempLink string

					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&project.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Dataflow Unit 2021 Assignment", Description: "Public"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/project", pavlosPyramid.Org.Path.URL), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response project.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Projects = append(pavlosPyramid.Projects, &response)
						tempLink = response.Path.URL
					})

					ginkgo.It("should be able to create a folder under the project", func() {
						payload, err = json.Marshal(&folder.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Dataflow Unit 2021 Solution", Description: "Private"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/folder", tempLink), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response folder.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Folders = append(pavlosPyramid.Folders, &response)
					})
				})

				ginkgo.When("I make a second create request", func() {
					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&project.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Forth Valley IT Department", Description: "Private"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/project", pavlosPyramid.Org.Path.URL), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response project.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Projects = append(pavlosPyramid.Projects, &response)
					})
				})
			})

			ginkgo.Context("who wants to create folders", func() {

				var tempLink string
				ginkgo.When("I make a valid create request", func() {
					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&folder.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "Forth Valley Assets"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/folder", pavlosPyramid.Org.Path.URL), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response folder.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Folders = append(pavlosPyramid.Folders, &response)
						tempLink = response.Path.URL
					})

					ginkgo.It("should create a project inside the folder", func() {
						payload, err = json.Marshal(&project.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "2020 Risk Analysis Project", Description: "Private"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/project", tempLink), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response project.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Projects = append(pavlosPyramid.Projects, &response)
					})
				})

				ginkgo.When("I make a second create request for folder", func() {
					ginkgo.It("should return the url of the resource", func() {
						payload, err = json.Marshal(&folder.CreateRequest{
							BaseInfo:   models.BaseInfo{Label: "2020 Risk Analysis", Description: "Private"},
							UserAction: models.UserAction{Owner: PavlosUser},
						})
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodPost, fmt.Sprintf("%s/folder", tempLink), payload).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusCreated))

						err = json.Unmarshal(recorder.Body.Bytes(), &content)
						gomega.Expect(err).To(gomega.BeNil())

						recorder = NewAPIRequest(http.MethodGet, content.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))

						var response folder.GetResponse
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						gomega.Expect(err).To(gomega.BeNil())

						pavlosPyramid.Folders = append(pavlosPyramid.Folders, &response)
					})

					ginkgo.It("should be able to go back using the previous link", func() {
						// select a random folder from the pyramid
						recorder = NewAPIRequest(http.MethodGet, pavlosPyramid.Folders[0].Path.URL, []byte("")).GetRecorder(web.Engine)
						gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))
					})
				})
			})

		})
	})
}
