package client

import (
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Test Suite")
}

var _ = BeforeSuite(func() {
	// Block all HTTP requests
	httpmock.Activate()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

var _ = BeforeEach(func() {
	// Remove any mocks
	httpmock.Reset()
})
