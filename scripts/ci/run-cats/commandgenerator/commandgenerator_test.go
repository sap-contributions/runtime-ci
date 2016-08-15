package commandgenerator_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator/commandgeneratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commandgenerator", func() {
	var nodes int
	var env *commandgeneratorfakes.FakeEnvironment

	BeforeEach(func() {
		env = &commandgeneratorfakes.FakeEnvironment{}
		nodes = 10
		env.GetNodesReturns(nodes, nil)
		env.GetSkipBackendCompatibilityReturns("backend_compatibility", nil)
		env.GetSkipDiegoDockerReturns("docker", nil)
		env.GetSkipInternetDependentReturns("internet_dependent", nil)
		env.GetSkipRouteServicesReturns("route_services", nil)
		env.GetSkipSecurityGroupsReturns("security_groups", nil)
		env.GetSkipServicesReturns("services", nil)
		env.GetSkipDiegoSSHReturns("ssh", nil)
		env.GetSkipV3Returns("v3", nil)
	})

	Context("When the path to CATs is set", func() {
		BeforeEach(func() {
			env.GetCatsPathReturns(".")
		})

		It("Should generate a command to run CATS", func() {
			cmd, args, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("bin/test"))

			Expect(strings.Join(args, " ")).To(Equal(
				fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=%d -skipPackage=backend_compatibility,docker,helpers,internet_dependent,route_services,security_groups,services,ssh,v3 -skip= -keepGoing", nodes),
			))

			env.GetCatsPathReturns("/path/to/cats")
			cmd, _, err = commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/path/to/cats/bin/test"))
		})

		Context("when the node count is unset", func() {
			BeforeEach(func() {
				env.GetNodesReturns(0, nil)
			})
			It("sets the default node count", func() {
				_, args, _ := commandgenerator.GenerateCmd(env)
				Expect(args).To(ContainElement("-nodes=2"))
			})
		})

		Context("when there are optional skipPackage env vars set", func() {
			Context("to not skip", func() {
				BeforeEach(func() {
					env.GetSkipBackendCompatibilityReturns("", nil)
					env.GetSkipDiegoDockerReturns("", nil)
					env.GetSkipInternetDependentReturns("", nil)
					env.GetSkipRouteServicesReturns("", nil)
					env.GetSkipSecurityGroupsReturns("", nil)
					env.GetSkipServicesReturns("", nil)
					env.GetSkipDiegoSSHReturns("", nil)
					env.GetSkipV3Returns("", nil)
				})

				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(" -skipPackage=helpers "))
				})
			})

			Context("to skip", func() {
				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(
						" -skipPackage=backend_compatibility,docker,helpers,internet_dependent,route_services,security_groups,services,ssh,v3 ",
					))
				})
			})

		})

		Context("when there are optional skip env vars set", func() {
			Context("to skip SSO", func() {
				BeforeEach(func() {
					env.GetSkipSSOReturns("SSO", nil)
				})

				It("generates a command that skips tests with the given tag", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement(ContainSubstring("-skip=SSO")))
				})
			})

			Context("to not skip SSO", func() {
				BeforeEach(func() {
					env.GetSkipSSOReturns("", nil)
				})

				It("generates a command that does not include the given tag in the skips", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).ToNot(ContainElement(ContainSubstring("-skip=SSO")))
				})
			})

		})
	})

	Context("When the path to CATS isn't explicitly provided", func() {
		BeforeEach(func() {
			env.GetGoPathReturns("/go")
		})

		It("Should return a sane default command path for use in Concourse", func() {
			cmd, _, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/go/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"))
		})
	})
})
