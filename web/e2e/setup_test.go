package e2e

import (
	"github.com/OBASHITechnology/resourceList/web"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestResourcePath(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "complete e2e Suite", []ginkgo.Reporter{reporters.NewJUnitReporter("./test-output.xml")})
}

var _ = ginkgo.BeforeSuite(func() {
	web.Registration()
})

var _ = ginkgo.Describe("Path using alias and Json content with URL prototype", func() {
	createResources()
})
