package bitcoin_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBitcoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitcoin Suite")
}
