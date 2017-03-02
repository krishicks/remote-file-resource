package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Suite")
}

var checkPath string
var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := gexec.Build("github.com/krishicks/remote-file-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())

	return []byte(path)
}, func(data []byte) {
	checkPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})
