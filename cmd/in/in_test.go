package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/krishicks/remote-file-resource/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("In", func() {
	var (
		session    *gexec.Session
		stdinBytes []byte
		cmd        *exec.Cmd
		targetPath string
	)

	BeforeEach(func() {
		var err error
		targetPath, err = ioutil.TempDir("", "in-test")
		Expect(err).NotTo(HaveOccurred())
		cmd = exec.Command(inPath, targetPath)
	})

	AfterEach(func() {
		os.RemoveAll(targetPath)
	})

	JustBeforeEach(func() {
		cmd.Stdin = bytes.NewReader(stdinBytes)

		var err error
		session, err = gexec.Start(
			cmd,
			GinkgoWriter,
			GinkgoWriter,
		)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when given invalid JSON", func() {
		BeforeEach(func() {
			stdinBytes = []byte("`")
		})

		It("logs an error to stderr", func() {
			Eventually(session.Err).Should(gbytes.Say("error reading request from stdin"))
		})
	})

	Context("when given valid JSON", func() {
		var (
			server *ghttp.Server
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.VerifyRequest("GET", "/releases/latest"),
			)

			request := types.InRequest{
				Source: types.Source{
					URI: server.URL() + "/releases/latest",
				},
				Version: types.Version{
					ETag: "expected-version",
				},
				Params: types.InParams{
					Filename: "release.tgz",
				},
			}

			var err error
			stdinBytes, err = json.Marshal(request)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			server.Close()
		})

		It("makes a request to the server", func() {
			Eventually(server.ReceivedRequests).Should(HaveLen(1))
		})

		Context("when the ETag of the artifact is the same as the version provided", func() {
			BeforeEach(func() {
				server.SetHandler(0,
					ghttp.RespondWith(
						http.StatusOK,
						"body",
						http.Header{
							"ETag": []string{"expected-version"},
						},
					),
				)
			})

			It("puts the resulting artifact in the target path", func() {
				path := filepath.Join(targetPath, "release.tgz")
				Eventually(func() bool {
					_, err := os.Lstat(path)
					return err == nil
				}).Should(BeTrue())

				bs, err := ioutil.ReadFile(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(bs).To(Equal([]byte("body")))
			})

			It("responds with the version that was downloaded", func() {
				expectedOutput := types.InResponse{
					Version: types.Version{
						ETag: "expected-version",
					},
				}

				bs, err := json.Marshal(expectedOutput)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session.Out).Should(gbytes.Say(string(bs)))
			})

			It("exits with no error", func() {
				<-session.Exited
				Expect(session.ExitCode()).To(Equal(0))
			})
		})

		Context("when the ETag of the artifact differs from the version provided", func() {
			BeforeEach(func() {
				server.SetHandler(0,
					ghttp.RespondWith(
						http.StatusOK,
						"body",
						http.Header{
							"ETag": []string{"unexpected-version"},
						},
					),
				)
			})

			It("logs an error to stderr", func() {
				Eventually(session.Err).Should(gbytes.Say("error downloading artifact; version expected-version is no longer available"))
			})

			It("exits with error", func() {
				<-session.Exited
				Expect(session.ExitCode()).To(Equal(1))
			})
		})
	})
})
