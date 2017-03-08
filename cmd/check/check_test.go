package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/krishicks/remote-file-resource/types"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check", func() {
	var (
		session    *gexec.Session
		stdinBytes []byte
		cmd        *exec.Cmd
	)

	BeforeEach(func() {
		cmd = exec.Command(checkPath)
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

		It("exits with error", func() {
			<-session.Exited
			Expect(session.ExitCode()).To(Equal(1))
		})
	})

	Context("when given valid JSON", func() {
		var (
			server *ghttp.Server
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.VerifyRequest("HEAD", "/releases/latest"),
			)

			request := types.CheckRequest{
				Source: types.Source{
					URI: server.URL() + "/releases/latest",
				},
				Version: types.Version{
					ETag: "old-version",
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
						"",
						http.Header{
							"ETag": []string{"old-version"},
						},
					),
				)
			})

			It("responds with the same version as was provided", func() {
				expectedOutput := types.CheckResponse{
					{ETag: "old-version"},
				}

				bs, err := json.Marshal(expectedOutput)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session.Out).Should(gbytes.Say(string(bs)))
			})
		})

		Context("when the ETag of the artifact differs from the version provided", func() {
			BeforeEach(func() {
				server.SetHandler(0,
					ghttp.RespondWith(
						http.StatusOK,
						"",
						http.Header{
							"ETag": []string{"new-version"},
						},
					),
				)
			})

			It("gets the artifact at the configured URI", func() {
				Eventually(server.ReceivedRequests).Should(HaveLen(1))
			})

			It("responds with the same version as was provided and the new version", func() {
				expectedOutput := types.CheckResponse{
					{ETag: "old-version"},
					{ETag: "new-version"},
				}

				bs, err := json.Marshal(expectedOutput)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session.Out.Contents).Should(MatchJSON(string(bs)))
			})
		})
	})
})
