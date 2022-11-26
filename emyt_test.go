package emyt_test

import (
	"github.com/emyt-io/emyt/dbprovider"
	"github.com/emyt-io/emyt/manager"
	. "github.com/onsi/ginkgo/v2"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var _ = Describe("Emyt", func() {

	var (
		t GinkgoTInterface
	)

	BeforeEach(func() {
		// Load dbprovider
		dbprovider.Init()
		t = GinkgoT()
	})

	Context("Successful Users", func() {
		It("Users should be set correctly", func() {
			apitest.New().
				Handler(manager.NewApp()).
				Get("/users").
				Expect(t).
				Status(http.StatusOK).
				Assert(func(res *http.Response, req *http.Request) error {
					assert.Equal(t, http.StatusOK, res.StatusCode)
					return nil
				}).
				End()
		})
	})

})
