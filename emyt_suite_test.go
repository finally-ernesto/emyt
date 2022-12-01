package emyt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEmyt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Emyt Suite")
}
