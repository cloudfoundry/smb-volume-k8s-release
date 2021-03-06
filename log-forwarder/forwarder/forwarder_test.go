package forwarder_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	. "code.cloudfoundry.org/volume-services-log-forwarder/forwarder"
	"code.cloudfoundry.org/volume-services-log-forwarder/forwarder/fluentshims/fluent_fake"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fluentshims/fluent_fake/fake_fluent.go ./fluentshims FluentInterface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./forwarder_fakes/fake_forwarder.go . Forwarder

var _ = Describe("Forwarder", func() {

	Describe("#Forward", func() {
		var (
			forwarder Forwarder
			err error
			fakeFluent *fluent_fake.FakeFluentInterface
			logger *lagertest.TestLogger
			appId string
			instanceId string
			log string
		)

		BeforeEach(func() {
			fakeFluent = &fluent_fake.FakeFluentInterface{}
			logger = lagertest.NewTestLogger("forwarder")
		})

		JustBeforeEach(func() {
			forwarder = NewForwarder(logger, fakeFluent)
			err = forwarder.Forward(appId, instanceId, log)
		})

		Context("given an event", func() {

			BeforeEach(func() {
				appId = "app-id"
				instanceId = "instance-id"
				log = "this is a test"
			})

			It("should post to fluentd", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeFluent.PostCallCount()).To(Equal(1))
				tag, message := fakeFluent.PostArgsForCall(0)
				Expect(tag).To(Equal("fluentd_dest"))
				Expect(message.(map[string]string)["app_id"]).To(Equal("app-id"))
				Expect(message.(map[string]string)["instance_id"]).To(Equal("instance-id"))
				Expect(message.(map[string]string)["source_type"]).To(Equal("VOL"))
				Expect(message.(map[string]string)["log"]).To(Equal("this is a test"))
			})

			It("should log", func(){
				Eventually(logger.Buffer()).Should(gbytes.Say("this is a test"))
			})
		})

		Context("when posting returns an error", func() {

			BeforeEach(func() {
				fakeFluent.PostReturns(errors.New("post-error"))
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
