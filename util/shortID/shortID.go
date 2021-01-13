package shortID

import (
	"github.com/lithammer/shortuuid"
)

// NewWithName accepts a name variable of string value in the following format
//  eg it will take the {uuid}/{uuid} -> convert it to {short uuid}
//  which will append it to the resource path URL. The above format define an internal path of a resource.
//  One example can be {org ID} / {repository ID} or {org ID} / {project ID} etc.
func NewWithURL(name string) string {
	return shortuuid.NewWithNamespace(name)
}
