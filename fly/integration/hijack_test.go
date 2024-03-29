package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/concourse/concourse/atc"
	"github.com/gorilla/websocket"
	"github.com/mgutz/ansi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Hijacking", func() {
	var hijacked <-chan struct{}
	var workingDirectory string
	var user string
	var path string
	var args []string

	BeforeEach(func() {
		hijacked = nil
		workingDirectory = ""
		user = "root"
		path = "bash"
		args = nil
	})

	upgrader := websocket.Upgrader{}

	hijackHandler := func(id string, didHijack chan<- struct{}, errorMessages []string, teamName string) http.HandlerFunc {
		return ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v1/teams/"+teamName+"/containers/%s/hijack", id)),
			func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()

				conn, err := upgrader.Upgrade(w, r, nil)
				Expect(err).NotTo(HaveOccurred())

				defer conn.Close()

				close(didHijack)

				var processSpec atc.HijackProcessSpec
				err = conn.ReadJSON(&processSpec)
				Expect(err).NotTo(HaveOccurred())

				Expect(processSpec.User).To(Equal(user))
				Expect(processSpec.Dir).To(Equal(workingDirectory))
				Expect(processSpec.Path).To(Equal(path))
				Expect(processSpec.Args).To(Equal(args))

				var payload atc.HijackInput
				err = conn.ReadJSON(&payload)
				Expect(err).NotTo(HaveOccurred())
				Expect(payload).To(Equal(atc.HijackInput{
					Stdin: []byte("some stdin"),
				}))

				err = conn.WriteJSON(atc.HijackOutput{
					Stdout: []byte("some stdout"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = conn.WriteJSON(atc.HijackOutput{
					Stderr: []byte("some stderr"),
				})
				Expect(err).NotTo(HaveOccurred())

				if len(errorMessages) > 0 {
					for _, msg := range errorMessages {
						err := conn.WriteJSON(atc.HijackOutput{
							Error: msg,
						})
						Expect(err).NotTo(HaveOccurred())
					}

					return
				}

				var closePayload atc.HijackInput
				err = conn.ReadJSON(&closePayload)
				Expect(err).NotTo(HaveOccurred())
				Expect(closePayload).To(Equal(atc.HijackInput{
					Closed: true,
				}))

				exitStatus := 123
				err = conn.WriteJSON(atc.HijackOutput{
					ExitStatus: &exitStatus,
				})
				Expect(err).NotTo(HaveOccurred())
			},
		)
	}

	fly := func(command string, args ...string) {
		commandWithArgs := append([]string{command}, args...)

		flyCmd := exec.Command(flyPath, append([]string{"-t", targetName}, commandWithArgs...)...)

		stdin, err := flyCmd.StdinPipe()
		Expect(err).NotTo(HaveOccurred())

		sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(hijacked).Should(BeClosed())

		_, err = fmt.Fprintf(stdin, "some stdin")
		Expect(err).NotTo(HaveOccurred())

		Eventually(sess.Out).Should(gbytes.Say("some stdout"))
		Eventually(sess.Err).Should(gbytes.Say("some stderr"))

		err = stdin.Close()
		Expect(err).NotTo(HaveOccurred())

		<-sess.Exited
		Expect(sess.ExitCode()).To(Equal(123))
	}

	hijack := func(args ...string) {
		fly("hijack", args...)
	}

	Context("with only a step name specified", func() {
		BeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/builds"),
					ghttp.RespondWithJSONEncoded(200, []atc.Build{
						{ID: 4, Name: "1", Status: "started", JobName: "some-job"},
						{ID: 3, Name: "3", Status: "started"},
						{ID: 2, Name: "2", Status: "started"},
						{ID: 1, Name: "1", Status: "finished"},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "build_id=3&step_name=some-step"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{
						{ID: "container-id-1", State: atc.ContainerStateCreated, BuildID: 3, Type: "task", StepName: "some-step", User: user},
					}),
				),
				hijackHandler("container-id-1", didHijack, nil, "main"),
			)
		})

		It("hijacks the most recent one-off build", func() {
			hijack("-s", "some-step")
		})

		It("hijacks the most recent one-off build with a more politically correct command", func() {
			fly("intercept", "-s", "some-step")
		})
	})

	Context("when the container specifies a working directory", func() {
		BeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack
			workingDirectory = "/tmp/build/my-favorite-guid"

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/builds"),
					ghttp.RespondWithJSONEncoded(200, []atc.Build{
						{ID: 3, Name: "3", Status: "started"},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "build_id=3&step_name=some-step"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{
						{ID: "container-id-1", State: atc.ContainerStateCreated, BuildID: 3, Type: "task", StepName: "some-step", WorkingDirectory: workingDirectory, User: user},
					}),
				),
				hijackHandler("container-id-1", didHijack, nil, "main"),
			)
		})

		It("hijacks the most recent one-off build in the specified working directory", func() {
			hijack("-s", "some-step")
		})
	})

	Context("when the container specifies a user", func() {
		BeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack
			user = "amelia"

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/builds"),
					ghttp.RespondWithJSONEncoded(200, []atc.Build{
						{ID: 3, Name: "3", Status: "started"},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "build_id=3&step_name=some-step"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{
						{ID: "container-id-1", State: atc.ContainerStateCreated, BuildID: 3, Type: "task", StepName: "some-step", User: "amelia"},
					}),
				),
				hijackHandler("container-id-1", didHijack, nil, "main"),
			)
		})

		It("hijacks the most recent one-off build as the specified user", func() {
			hijack("-s", "some-step")
		})
	})

	Context("when no containers are found", func() {
		BeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/builds"),
					ghttp.RespondWithJSONEncoded(200, []atc.Build{
						{ID: 1, Name: "1", Status: "finished"},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "build_id=1&step_name=some-step"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{}),
				),
				hijackHandler("container-id-1", didHijack, nil, "main"),
			)
		})

		It("return a friendly error message", func() {
			flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-s", "some-step")
			sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(1))

			Expect(sess.Err).To(gbytes.Say("no containers matched your search parameters!\n\nthey may have expired if your build hasn't recently finished.\n"))
		})

		Context("when a url is passed", func() {
			It("return a friendly error message", func() {
				flyCmd := exec.Command(flyPath, "hijack", "-s", "some-step", "-u", fmt.Sprintf("%s/teams/%s", atcServer.URL(), teamName))
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(1))

				Expect(sess.Err).To(gbytes.Say("no containers matched your search parameters!\n\nthey may have expired if your build hasn't recently finished.\n"))
			})

			It("returns an error when target from url is not found", func() {
				flyCmd := exec.Command(flyPath, "hijack", "-s", "some-step", "-u", fmt.Sprintf("%s/teams/%s", "http://faketarget.com", teamName))
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(1))

				Expect(sess.Err).To(gbytes.Say("no target matching url"))
			})

			It("returns an error when team name from url is not found", func() {
				flyCmd := exec.Command(flyPath, "hijack", "-s", "some-step", "-u", fmt.Sprintf("%s/teams/%s/builds/0", atcServer.URL(), "faketeam"))
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(1))

				Expect(sess.Err).To(gbytes.Say("no target matching url"))
			})
		})
	})

	Context("when no containers are found", func() {
		BeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack
			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "build_id=0"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{}),
				),
			)
		})

		It("logs an error message and response status/body", func() {
			flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-b", "0")

			stdin, err := flyCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			defer stdin.Close()

			sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess.Err.Contents).Should(ContainSubstring("no containers matched your search parameters!\n\nthey may have expired if your build hasn't recently finished.\n"))

			<-sess.Exited
			Expect(sess.ExitCode()).To(Equal(1))
		})
	})

	Context("when multiple step containers are found", func() {
		var (
			containerList       []atc.Container
			didHijack           chan struct{}
			expectedQueryParams []string
		)

		BeforeEach(func() {
			didHijack = make(chan struct{})
			hijacked = didHijack
			containerList = []atc.Container{
				{
					ID:           "container-id-1",
					WorkerName:   "worker-name-1",
					PipelineName: "pipeline-name-1",
					JobName:      "some-job",
					BuildName:    "2",
					BuildID:      12,
					Type:         "get",
					StepName:     "some-input",
					Attempt:      "1.1.1",
					User:         user,
					State:        atc.ContainerStateCreated,
				},
				{
					ID:           "container-id-2",
					WorkerName:   "worker-name-2",
					PipelineName: "pipeline-name-1",
					JobName:      "some-job",
					BuildName:    "2",
					BuildID:      13,
					Type:         "put",
					StepName:     "some-output",
					Attempt:      "1.1.2",
					User:         user,
					State:        atc.ContainerStateCreated,
				},
				{
					ID:           "container-id-3",
					WorkerName:   "worker-name-2",
					PipelineName: "pipeline-name-2",
					JobName:      "some-job",
					BuildName:    "2",
					BuildID:      13,
					StepName:     "some-output",
					Type:         "task",
					Attempt:      "1",
					User:         user,
					State:        atc.ContainerStateCreated,
				},
				{
					ID:           "container-id-4",
					WorkerName:   "worker-name-2",
					PipelineName: "pipeline-name-2",
					ResourceName: "banana",
					User:         user,
					Type:         "check",
					State:        atc.ContainerStateCreated,
				},
			}
			expectedQueryParams = append([]string{}, "pipeline_name=pipeline-name-1", "job_name=some-job")
		})

		JustBeforeEach(func() {
			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", strings.Join(expectedQueryParams, "&")),
					ghttp.RespondWithJSONEncoded(200, containerList),
				),
				hijackHandler("container-id-2", didHijack, nil, "main"),
			)
		})

		It("asks the user to select the container from a menu", func() {
			flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-j", "pipeline-name-1/some-job")

			stdin, err := flyCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess.Out).Should(gbytes.Say("1. resource: banana, type: check"))
			Eventually(sess.Out).Should(gbytes.Say("2. build #2, step: some-input, type: get, attempt: 1.1.1"))
			Eventually(sess.Out).Should(gbytes.Say("3. build #2, step: some-output, type: put, attempt: 1.1.2"))
			Eventually(sess.Out).Should(gbytes.Say("4. build #2, step: some-output, type: task, attempt: 1"))
			Eventually(sess.Out).Should(gbytes.Say("choose a container: "))

			_, err = fmt.Fprintf(stdin, "3\n")
			Expect(err).NotTo(HaveOccurred())

			Eventually(hijacked).Should(BeClosed())

			_, err = fmt.Fprintf(stdin, "some stdin")
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess.Out).Should(gbytes.Say("some stdout"))
			Eventually(sess.Err).Should(gbytes.Say("some stderr"))

			err = stdin.Close()
			Expect(err).NotTo(HaveOccurred())

			<-sess.Exited
			Expect(sess.ExitCode()).To(Equal(123))
		})

		Context("and no containers are in hijackable state", func() {
			BeforeEach(func() {
				containerList = []atc.Container{
					{
						ID:           "container-id-2",
						WorkerName:   "worker-name-1",
						PipelineName: "pipeline-name-1",
						JobName:      "some-job",
						BuildName:    "2",
						BuildID:      12,
						Type:         "get",
						StepName:     "some-input",
						Attempt:      "1.1.1",
						User:         user,
						State:        atc.ContainerStateCreating,
					},
				}
			})

			It("should show that no containers are hijackable", func() {
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-j", "pipeline-name-1/some-job")

				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(1))

				Eventually(sess.Err).Should(gbytes.Say("no containers matched"))
				close(didHijack)
			})
		})

		Context("and some containers are in a non-hijackable state", func() {
			BeforeEach(func() {
				containerList = []atc.Container{
					{
						ID:           "container-id-1",
						WorkerName:   "worker-name-1",
						PipelineName: "pipeline-name-1",
						JobName:      "some-job",
						BuildName:    "2",
						BuildID:      12,
						Type:         "get",
						StepName:     "some-input",
						Attempt:      "1.1.1",
						User:         user,
						State:        atc.ContainerStateCreating,
					},
					{
						ID:           "container-id-2",
						WorkerName:   "worker-name-2",
						PipelineName: "pipeline-name-1",
						JobName:      "some-job",
						BuildName:    "2",
						BuildID:      13,
						Type:         "put",
						StepName:     "some-output",
						Attempt:      "1.1.2",
						User:         user,
						State:        atc.ContainerStateCreated,
					},
					{
						ID:           "container-id-3",
						WorkerName:   "worker-name-2",
						PipelineName: "pipeline-name-2",
						JobName:      "some-job",
						BuildName:    "2",
						BuildID:      13,
						StepName:     "some-output",
						Type:         "task",
						Attempt:      "1",
						User:         user,
						State:        atc.ContainerStateFailed,
					},
					{
						ID:           "container-id-4",
						WorkerName:   "worker-name-2",
						PipelineName: "pipeline-name-2",
						ResourceName: "banana",
						User:         user,
						Type:         "check",
						State:        atc.ContainerStateDestroying,
					},
				}
			})

			It("should not display those containers in the list of results", func() {
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-j", "pipeline-name-1/some-job")

				stdin, err := flyCmd.StdinPipe()
				Expect(err).NotTo(HaveOccurred())

				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Out).Should(gbytes.Say("1. build #2, step: some-output, type: put, attempt: 1.1.2"))
				Eventually(sess.Out).Should(gbytes.Say("2. build #2, step: some-output, type: task, attempt: 1"))
				Eventually(sess.Out).Should(gbytes.Say("choose a container: "))

				_, err = fmt.Fprintf(stdin, "1\n")
				Expect(err).NotTo(HaveOccurred())

				Eventually(hijacked).Should(BeClosed())

				_, err = fmt.Fprintf(stdin, "some stdin")
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Out).Should(gbytes.Say("some stdout"))
				Eventually(sess.Err).Should(gbytes.Say("some stderr"))

				err = stdin.Close()
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(123))
			})
		})

		Context("and only one container is in hijackable state", func() {
			BeforeEach(func() {
				containerList = []atc.Container{
					{
						ID:           "container-id-1",
						WorkerName:   "worker-name-1",
						PipelineName: "pipeline-name-1",
						JobName:      "some-job",
						BuildName:    "1",
						BuildID:      12,
						Type:         "get",
						StepName:     "some-input",
						Attempt:      "1.1.1",
						User:         user,
						State:        atc.ContainerStateDestroying,
					},
					{
						ID:           "container-id-2",
						WorkerName:   "worker-name-2",
						PipelineName: "pipeline-name-1",
						JobName:      "some-job",
						BuildName:    "2",
						BuildID:      13,
						Type:         "put",
						StepName:     "some-output",
						Attempt:      "1.1.2",
						User:         user,
						State:        atc.ContainerStateCreated,
					},
				}
			})

			It("hijacks the hijackable container", func() {
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "-j", "pipeline-name-1/some-job")

				stdin, err := flyCmd.StdinPipe()
				Expect(err).NotTo(HaveOccurred())

				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(hijacked).Should(BeClosed())

				_, err = fmt.Fprintf(stdin, "some stdin")
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Out).Should(gbytes.Say("some stdout"))
				Eventually(sess.Err).Should(gbytes.Say("some stderr"))

				err = stdin.Close()
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(123))
			})
		})
	})

	Context("when hijack returns a single container", func() {
		var (
			containerArguments []string
			stepType           string
			stepName           string
			buildID            int
			hijackHandlerError []string
			pipelineName       string
			resourceName       string
			jobName            string
			buildName          string
			attempt            string
			hijackTeamName     string
		)

		BeforeEach(func() {
			hijackHandlerError = nil
			pipelineName = "a-pipeline"
			jobName = ""
			buildName = ""
			buildID = 0
			stepType = ""
			stepName = ""
			resourceName = ""
			containerArguments = []string{}
			attempt = ""
			hijackTeamName = "main"
		})

		JustBeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName+"/containers", strings.Join(containerArguments, "&")),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{
						{ID: "container-id-1", State: atc.ContainerStateCreated, WorkerName: "some-worker", PipelineName: pipelineName, JobName: jobName, BuildName: buildName, BuildID: buildID, Type: stepType, StepName: stepName, ResourceName: resourceName, Attempt: attempt, User: user},
					}),
				),
				hijackHandler("container-id-1", didHijack, hijackHandlerError, hijackTeamName),
			)
		})

		Context("when called with check container", func() {
			BeforeEach(func() {
				resourceName = "some-resource-name"
				containerArguments = append(containerArguments, "type=check", "resource_name=some-resource-name", "pipeline_name=a-pipeline")
			})

			Context("when the team is 'main'", func() {
				BeforeEach(func() {
					hijackTeamName = "main"
				})
				Context("and with pipeline specified", func() {
					It("can accept the check resources name and a pipeline", func() {
						hijack("--check", "a-pipeline/some-resource-name")
					})

					It("hijacks the given check container by URL", func() {
						hijack("--url", atcServer.URL()+"/teams/"+teamName+"/pipelines/a-pipeline/resources/some-resource-name")
					})

					Context("and with pipeline instance is specified", func() {
						BeforeEach(func() {
							containerArguments = append(containerArguments, "vars=%7B%22branch%22%3A%22master%22%7D")
						})

						It("can accept the check resources name and a pipeline", func() {
							hijack("--check", "a-pipeline/branch:master/some-resource-name")
						})

						It("hijacks the given check container by URL", func() {
							hijack("--url", atcServer.URL()+"/teams/"+teamName+"/pipelines/a-pipeline/resources/some-resource-name"+"?vars.branch=%22master%22")
						})
					})
				})
			})

			Context("when the team is 'other'", func() {
				BeforeEach(func() {
					hijackTeamName = "other"

					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
							ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
					)
				})

				Context("and with pipeline specified", func() {
					It("can accept the check resources name and a pipeline", func() {
						hijack("--check", "a-pipeline/some-resource-name", "--team", hijackTeamName)
					})
				})

				Context("and with url specified", func() {
					It("hijacks the given check container by URL", func() {
						hijack("--url", atcServer.URL()+"/teams/other/pipelines/a-pipeline/resources/some-resource-name", "--team", "other")
					})
				})
			})
		})

		Context("when called with a specific build id", func() {
			BeforeEach(func() {
				containerArguments = append(containerArguments, "build_id=2", "step_name=some-step")
				stepType = "task"
				stepName = "some-step"
				buildID = 2
			})

			Context("when the team is 'main'", func() {
				BeforeEach(func() {
					hijackTeamName = "main"
				})
				It("hijacks the most recent one-off build", func() {
					hijack("-b", "2", "-s", "some-step")
				})
			})

			Context("when the team is 'other'", func() {
				BeforeEach(func() {
					hijackTeamName = "other"

					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
							ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
					)
				})

				It("hijacks the most recent one-off build", func() {
					hijack("-b", "2", "-s", "some-step", "--team", hijackTeamName)
				})
			})
		})

		Context("when called with a specific job", func() {
			BeforeEach(func() {
				containerArguments = append(containerArguments, "pipeline_name=some-pipeline", "job_name=some-job", "step_name=some-step")
				jobName = "some-job"
				buildName = "3"
				buildID = 13
				stepType = "task"
				stepName = "some-step"
			})

			Context("hijacks the job's next build", func() {
				Context("When the team is 'main'", func() {
					BeforeEach(func() {
						hijackTeamName = "main"
					})

					It("hijacks the job's next build with '<pipeline>/<job>'", func() {
						hijack("--job", "some-pipeline/some-job", "--step", "some-step")
					})

					It("hijacks the job's next build when URL is specified", func() {
						hijack("--url", atcServer.URL()+"/teams/"+teamName+"/pipelines/some-pipeline/jobs/some-job", "--step", "some-step")
					})

					Context("when pipeline instance is specified", func() {
						BeforeEach(func() {
							containerArguments = append(containerArguments, "vars=%7B%22branch%22%3A%22master%22%2C%22other.field%22%3A123%7D")
						})

						It("hijacks the job's next build with '<pipeline>/<instance_vars>/<job>'", func() {
							hijack("--job", `some-pipeline/branch:master,"other.field":123/some-job`, "--step", "some-step")
						})

						It("hijacks the job's next build when URL is specified", func() {
							hijack("--url", atcServer.URL()+"/teams/"+teamName+"/pipelines/some-pipeline/jobs/some-job"+`?vars.branch="master"&vars."other.field"=123&ignore=true`, "--step", "some-step")
						})
					})

				})

				Context("When the team is 'other'", func() {
					BeforeEach(func() {
						hijackTeamName = "other"

						atcServer.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
								ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
						)
					})
					It("hijacks the job's next build with 'pipelineName/jobName'", func() {
						hijack("--job", "some-pipeline/some-job", "--step", "some-step", "--team", hijackTeamName)
					})

					It("hijacks the job's next build when URL is specified", func() {
						hijack("--url", atcServer.URL()+"/teams/other/pipelines/some-pipeline/jobs/some-job", "--step", "some-step", "--team", hijackTeamName)
					})
				})

			})

			Context("with a specific build of the job", func() {
				BeforeEach(func() {
					containerArguments = append(containerArguments, "build_name=3")
				})

				Context("When the team is 'main'", func() {
					BeforeEach(func() {
						hijackTeamName = "main"
					})
					It("hijacks the given build", func() {
						hijack("--job", "some-pipeline/some-job", "--build", "3", "--step", "some-step")
					})
					It("hijacks the given build with URL", func() {
						hijack("--url", atcServer.URL()+"/teams/main/pipelines/some-pipeline/jobs/some-job/builds/3", "--step", "some-step")
					})
				})

				Context("When the team is 'other'", func() {
					BeforeEach(func() {
						hijackTeamName = "other"

						atcServer.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
								ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
						)
					})
					It("hijacks the given build", func() {
						hijack("--job", "some-pipeline/some-job", "--build", "3", "--step", "some-step", "--team", hijackTeamName)
					})
					It("hijacks the given build with URL", func() {
						hijack("--url", atcServer.URL()+"/teams/other/pipelines/some-pipeline/jobs/some-job/builds/3", "--step", "some-step", "--team", hijackTeamName)
					})
				})

			})
		})

		Context("when called with a specific attempt number", func() {
			BeforeEach(func() {
				containerArguments = append(containerArguments, "pipeline_name=some-pipeline", "job_name=some-job", "step_name=some-step", "attempt=2.4")
				jobName = "some-job"
				buildName = "3"
				buildID = 13
				stepType = "task"
				stepName = "some-step"
				attempt = "2.4"
			})

			Context("When the team is 'main'", func() {
				BeforeEach(func() {
					hijackTeamName = "main"
				})
				It("hijacks the job's next build", func() {
					hijack("--job", "some-pipeline/some-job", "--step", "some-step", "--attempt", "2.4")
				})
			})

			Context("When the team is 'other'", func() {
				BeforeEach(func() {
					hijackTeamName = "other"

					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
							ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
					)
				})
				It("hijacks the job's next build", func() {
					hijack("--job", "some-pipeline/some-job", "--step", "some-step", "--attempt", "2.4", "--team", hijackTeamName)
				})
			})
		})

		Context("when called with a step type", func() {
			BeforeEach(func() {
				containerArguments = append(containerArguments, "pipeline_name=some-pipeline", "job_name=some-job", "step_name=some-step", "type=put")
				jobName = "some-job"
				buildName = "3"
				buildID = 13
				stepType = "put"
				stepName = "some-step"
				attempt = ""
			})

			It("hijacks the job's next build", func() {
				hijack("--job", "some-pipeline/some-job", "--step", "some-step", "--step-type", "put")
			})
		})

		Context("when called with a specific path and args", func() {
			BeforeEach(func() {
				path = "sh"
				args = []string{"echo hello"}

				containerArguments = append(containerArguments, "build_id=2", "step_name=some-step")
				stepType = "task"
				stepName = "some-step"
				buildID = 2
			})

			It("hijacks and runs the provided path with args", func() {
				hijack("-b", "2", "-s", "some-step", "sh", "echo hello")
			})
		})

		Context("when hijacking yields an error", func() {
			BeforeEach(func() {
				resourceName = "some-resource-name"
				containerArguments = append(containerArguments, "type=check", "resource_name=some-resource-name", "pipeline_name=a-pipeline")
				hijackHandlerError = []string{"something went wrong"}
			})

			It("prints it to stderr and exits 255", func() {
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--check", "a-pipeline/some-resource-name")

				stdin, err := flyCmd.StdinPipe()
				Expect(err).NotTo(HaveOccurred())

				defer stdin.Close()

				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(hijacked).Should(BeClosed())

				_, err = fmt.Fprintf(stdin, "some stdin")
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Err.Contents).Should(ContainSubstring(ansi.Color("something went wrong", "red+b") + "\n"))

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(255))
			})
		})
	})

	Context("when hijacking a specific container", func() {
		var (
			hijackHandlerError []string
			statusCode         int
			id                 string
			hijackTeamName     string
		)

		BeforeEach(func() {
			hijackHandlerError = nil
			statusCode = 0
			id = ""
			hijackTeamName = "main"

		})

		JustBeforeEach(func() {
			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName+"/containers/container-id"),
					ghttp.RespondWithJSONEncoded(statusCode, atc.Container{
						ID:   id,
						User: user,
					}),
				),
			)
		})

		Context("when container exists", func() {
			BeforeEach(func() {
				statusCode = 200
				id = "container-id"
			})

			Context("when hijack returns no error", func() {
				JustBeforeEach(func() {
					didHijack := make(chan struct{})
					hijacked = didHijack
					atcServer.AppendHandlers(
						hijackHandler("container-id", didHijack, hijackHandlerError, hijackTeamName),
					)
				})
				Context("when the team is 'main'", func() {
					BeforeEach(func() {
						hijackTeamName = "main"
					})
					It("should hijack container with associated handle", func() {
						hijack("--handle", "container-id")
					})
				})

				Context("when the team is 'other'", func() {
					BeforeEach(func() {
						hijackTeamName = "other"

						atcServer.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/api/v1/teams/"+hijackTeamName),
								ghttp.RespondWithJSONEncoded(http.StatusOK, atc.Team{Name: hijackTeamName})),
						)
					})
					It("should hijack container with associated handle to 'other' team", func() {
						hijack("--handle", "container-id", "--team", hijackTeamName)
					})
				})
			})

			Context("when hijack returns error", func() {
				JustBeforeEach(func() {
					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", fmt.Sprintf("/api/v1/teams/main/containers/%s/hijack", id)),
							ghttp.RespondWithJSONEncoded(403, nil),
						),
					)
				})

				It("should print out response status and error", func() {
					flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--handle", "container-id")
					sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					Eventually(sess).Should(gexec.Exit(1))

					Expect(sess.Err).To(gbytes.Say("error: 403 Forbidden websocket: bad handshake"))
				})
			})
		})

		Context("when container does not exist", func() {
			BeforeEach(func() {
				statusCode = 404
			})

			JustBeforeEach(func() {
				didHijack := make(chan struct{})
				hijacked = didHijack
				atcServer.AppendHandlers(
					hijackHandler("container-id", didHijack, hijackHandlerError, hijackTeamName),
				)
			})

			It("should return an appropriate error message", func() {
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--handle", "container-id")
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(1))

				Expect(sess.Err).To(gbytes.Say("no containers matched the given handle id!\n\nthey may have expired if your build hasn't recently finished.\n"))
			})
		})
	})

	Context("when passing a URL that doesn't match the target", func() {
		It("errors out when wrong team is specified", func() {
			flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--url", atcServer.URL()+"/teams/wrongteam/pipelines/a-pipeline/resources/some-resource-name")

			sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess.Err.Contents).Should(ContainSubstring("Team in URL doesn't match the current team of the target"))

			<-sess.Exited
			Expect(sess.ExitCode()).ToNot(Equal(0))
		})

		It("errors out when wrong URL is specified", func() {
			flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--url", "http://wrong.example.com/teams/"+teamName+"/pipelines/a-pipeline/resources/some-resource-name")

			sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess.Err.Contents).Should(ContainSubstring("URL doesn't match that of target"))

			<-sess.Exited
			Expect(sess.ExitCode()).ToNot(Equal(0))
		})
	})

	Context("when hijacking yields an executable not found error", func() {
		var hijacked2 <-chan struct{}
		JustBeforeEach(func() {
			didHijack := make(chan struct{})
			hijacked = didHijack
			didHijack2 := make(chan struct{})
			hijacked2 = didHijack2
			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers", "type=check&resource_name=some-resource-name&pipeline_name=a-pipeline"),
					ghttp.RespondWithJSONEncoded(200, []atc.Container{
						{ID: "container-id-1", State: atc.ContainerStateCreated, WorkerName: "some-worker", PipelineName: "a-pipeline", JobName: "", BuildName: "", BuildID: 0, Type: "", StepName: "", ResourceName: "some-resource-name", Attempt: "", User: user},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/teams/main/containers/container-id-1/hijack"),
					func(w http.ResponseWriter, r *http.Request) {
						defer GinkgoRecover()

						conn, err := upgrader.Upgrade(w, r, nil)
						Expect(err).NotTo(HaveOccurred())

						defer conn.Close()

						close(didHijack)

						var processSpec atc.HijackProcessSpec
						err = conn.ReadJSON(&processSpec)
						Expect(err).NotTo(HaveOccurred())

						Expect(processSpec.User).To(Equal(user))
						Expect(processSpec.Dir).To(Equal(workingDirectory))
						Expect(processSpec.Path).To(Equal("bash"))
						Expect(processSpec.Args).To(Equal(args))

						err = conn.WriteJSON(atc.HijackOutput{
							ExecutableNotFound: true,
						})
						Expect(err).NotTo(HaveOccurred())

						err = conn.WriteJSON(atc.HijackOutput{
							Error: "executable not found",
						})
						Expect(err).NotTo(HaveOccurred())
					},
				),
				hijackHandler("container-id-1", didHijack2, nil, "main"),
			)
		})

		Context("when a path was not specified", func() {
			BeforeEach(func() {
				path = "sh"
			})
			It("tries \"bash\" then \"sh\"", func() {
				os.Stdout.WriteString("\n")
				flyCmd := exec.Command(flyPath, "-t", targetName, "hijack", "--check", "a-pipeline/some-resource-name")

				stdin, err := flyCmd.StdinPipe()
				Expect(err).NotTo(HaveOccurred())

				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(hijacked).Should(BeClosed())

				Eventually(sess.Err.Contents).Should(ContainSubstring(ansi.Color("executable not found", "red+b") + "\n"))
				Eventually(sess.Err.Contents).Should(ContainSubstring("Couldn't find \"bash\" on container, retrying with \"sh\""))

				Eventually(hijacked2).Should(BeClosed())

				_, err = fmt.Fprintf(stdin, "some stdin")
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Out).Should(gbytes.Say("some stdout"))
				Eventually(sess.Err).Should(gbytes.Say("some stderr"))

				err = stdin.Close()
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(123))
			})
		})
	})
})
