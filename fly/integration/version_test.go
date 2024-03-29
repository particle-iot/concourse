package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/concourse/concourse/atc"

	"github.com/concourse/concourse/fly/version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Version Checks", func() {
	// patch version
	var (
		flyVersion       string
		customAtcVersion string
		flySession       *gexec.Session
	)
	BeforeEach(func() {
		flyVersion = atcVersion

		atcServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers"),
				ghttp.RespondWith(http.StatusOK, "[]"),
			),
		)
	})

	JustBeforeEach(func() {
		atcServer.SetHandler(3,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/info"),
				ghttp.RespondWithJSONEncoded(200, atc.Info{Version: customAtcVersion, WorkerVersion: workerVersion}),
			),
		)

		flyCmd := exec.Command(flyPath, "-t", targetName, "containers")
		flyCmd.Env = append(os.Environ(), "FAKE_FLY_VERSION="+flyVersion)

		var err error
		flySession, err = gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("when the client and server differ by a patch version", func() {
		BeforeEach(func() {
			major, minor, patch, err := version.GetSemver(atcVersion)
			Expect(err).NotTo(HaveOccurred())

			customAtcVersion = fmt.Sprintf("%d.%d.%d", major, minor, patch+1)
		})

		It("warns the user that there is a difference", func() {
			Eventually(flySession).Should(gexec.Exit(0))
			Expect(flySession.Err).To(gbytes.Say(`fly version \(%s\) is out of sync with the target \(%s\). to sync up, run the following:\n\n    `, flyVersion, customAtcVersion))
			Expect(flySession.Err).To(gbytes.Say(`fly.* -t %s sync\n`, targetName))
		})
	})

	// when then match
	Describe("when the client and server are the same version", func() {
		BeforeEach(func() {
			customAtcVersion = atcVersion
		})

		It("it doesn't give any warning message", func() {
			Eventually(flySession).Should(gexec.Exit(0))
			Expect(flySession.Err).ShouldNot(gbytes.Say("version"))
		})
	})

	// minor version
	Describe("when the client and server differ by a minor version", func() {
		BeforeEach(func() {
			major, minor, patch, err := version.GetSemver(atcVersion)
			Expect(err).NotTo(HaveOccurred())

			customAtcVersion = fmt.Sprintf("%d.%d.%d", major, minor+1, patch)
		})

		It("error and tell the user to sync", func() {
			Eventually(flySession).Should(gexec.Exit(1))
			Expect(flySession.Err).To(gbytes.Say(`fly version \(%s\) is out of sync with the target \(%s\). to sync up, run the following:\n\n    `, flyVersion, customAtcVersion))
			Expect(flySession.Err).To(gbytes.Say(`fly.* -t %s sync\n`, targetName))
			Expect(flySession.Err).To(gbytes.Say("cowardly refusing to run due to significant version discrepancy"))
		})
	})

	// major version (same as minor)
	Describe("when the client and server differ by a major version", func() {
		BeforeEach(func() {
			major, minor, patch, err := version.GetSemver(atcVersion)
			Expect(err).NotTo(HaveOccurred())

			customAtcVersion = fmt.Sprintf("%d.%d.%d", major+1, minor, patch)
		})

		It("error and tell the user to sync", func() {
			Eventually(flySession).Should(gexec.Exit(1))
			Expect(flySession.Err).To(gbytes.Say(`fly version \(%s\) is out of sync with the target \(%s\). to sync up, run the following:\n\n    `, flyVersion, customAtcVersion))
			Expect(flySession.Err).To(gbytes.Say(`fly.* -t %s sync\n`, targetName))
			Expect(flySession.Err).To(gbytes.Say("cowardly refusing to run due to significant version discrepancy"))
		})
	})

	// dev version
	Describe("when the client is a development version", func() {
		BeforeEach(func() {
			flyVersion = "0.0.0-dev"
		})

		It("never complains", func() {
			Eventually(flySession).Should(gexec.Exit(0))
			Expect(flySession.Err).ShouldNot(gbytes.Say("version"))
		})
	})
})
