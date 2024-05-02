package cmd

import (
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var r *rand.Rand

func TestCmd(t *testing.T) {

	BeforeSuite(func() {

		r = rand.New(rand.NewSource(GinkgoRandomSeed()))
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
