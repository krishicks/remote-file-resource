package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}

var inPath string
var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := gexec.Build("github.com/krishicks/remote-file-resource/cmd/in")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	inPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})
