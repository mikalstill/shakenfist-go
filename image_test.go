package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image management functions", func() {
	const (
		test_url       string = "http://server:13000"
		test_namespace string = "testspace"
		test_key       string = "testkey"
	)

	var (
		client *Client
	)

	BeforeEach(func() {
		// Configure client
		client = NewClient(test_url, test_namespace, test_key)

		httpmock.RegisterResponder("POST", test_url+"/auth",
			httpmock.NewBytesResponder(200, []byte(`{"access_token":"ABC123"}`)))
	})

	It("should get an image metadata list", func() {
		testImageMeta := []byte(`[
			{
			  "checksum": "ed44b9745b8d62bcbbc180b5f36c24bb",
			  "fetched": "Wed, 21 Oct 2020 10:08:16 -0000",
			  "file_version": 1,
			  "modified": "Fri, 16 Oct 2020 16:32:30 GMT",
			  "size": "359464960",
			  "url": "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
			  "ref": "095fdd2b66627f1665a53623c77d00f82cd373602a0b445470ac0437885412aa",
			  "node": "sf-2"
			},
			{
			  "checksum": "ff26abf6a7b47feeeb34364bb915160d",
			  "fetched": "Fri, 23 Oct 2020 23:56:29 -0000",
			  "file_version": 2,
			  "modified": "Thu, 22 Oct 2020 15:11:56 GMT",
			  "size": "558760448",
			  "url": "https://cloud-images.ubuntu.com/groovy/current/groovy-server-cloudimg-amd64.img",
			  "ref": "1b01f4bcb02f3a060610a4f73b34012d59197a12c2794b495dd583e43d0f65e8",
			  "node": "sf-3"
			}
		  ]`)

		// Prepare mocked HTTP
		reqPath := test_url + "/images"
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testImageMeta))

		// Make client request
		image, err := client.GetImageMeta()
		Expect(err).To(BeNil())
		Expect(image).To(Equal([]ImageMeta{
			{
				Checksum:    "ed44b9745b8d62bcbbc180b5f36c24bb",
				Fetched:     "Wed, 21 Oct 2020 10:08:16 -0000",
				FileVersion: 1,
				Modified:    "Fri, 16 Oct 2020 16:32:30 GMT",
				Size:        "359464960",
				URL:         "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
				Ref:         "095fdd2b66627f1665a53623c77d00f82cd373602a0b445470ac0437885412aa",
				Node:        "sf-2",
			},
			{
				Checksum:    "ff26abf6a7b47feeeb34364bb915160d",
				Fetched:     "Fri, 23 Oct 2020 23:56:29 -0000",
				FileVersion: 2,
				Modified:    "Thu, 22 Oct 2020 15:11:56 GMT",
				Size:        "558760448",
				URL:         "https://cloud-images.ubuntu.com/groovy/current/groovy-server-cloudimg-amd64.img",
				Ref:         "1b01f4bcb02f3a060610a4f73b34012d59197a12c2794b495dd583e43d0f65e8",
				Node:        "sf-3",
			},
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})
})
