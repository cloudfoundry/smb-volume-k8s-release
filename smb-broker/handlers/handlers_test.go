package handlers_test

import (
	. "code.cloudfoundry.org/smb-broker/handlers"
	smbbrokerfakes "code.cloudfoundry.org/smb-broker/smb-brokerfakes"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6  k8s.io/client-go/kubernetes/typed/core/v1.PersistentVolumeInterface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6  k8s.io/client-go/kubernetes/typed/core/v1.PersistentVolumeClaimInterface

var _ = Describe("Handlers", func() {
	var brokerHandler http.Handler
	var err error
	var recorder *httptest.ResponseRecorder
	var request *http.Request
	var fakePersitentVolumeClient *smbbrokerfakes.FakePersistentVolumeInterface
	var fakePersitentVolumeClaimClient *smbbrokerfakes.FakePersistentVolumeClaimInterface

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		fakePersitentVolumeClient = &smbbrokerfakes.FakePersistentVolumeInterface{}
		fakePersitentVolumeClaimClient = &smbbrokerfakes.FakePersistentVolumeClaimInterface{}
	})

	JustBeforeEach(func() {
		brokerHandler, err = BrokerHandler(fakePersitentVolumeClient, fakePersitentVolumeClaimClient)
	})

	Describe("Endpoints", func() {
		var source = rand.NewSource(GinkgoRandomSeed())
		JustBeforeEach(func() {
			brokerHandler.ServeHTTP(recorder, request)
		})

		Describe("#Catalog endpoint", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/v2/catalog", nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should list catalog of services offered by the SMB service broker", func() {
				Expect(recorder.Code).To(Equal(200))
				Expect(recorder.Body).To(MatchJSON(fixture("catalog.json")))
			})
		})

		Describe("#Provision endpoint", func() {
			var serviceInstanceKey string

			BeforeEach(func() {
				serviceInstanceKey = randomString(source)

				var err error
				request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": { "parameter1": "1", "parameter2": "foo" } }`))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should allow provisioning", func() {
				Expect(recorder.Code).To(Equal(201))
				Expect(recorder.Body).To(MatchJSON(`{}`))
			})

			It("should create a persistent volume", func() {
				Expect(fakePersitentVolumeClient.CreateCallCount()).To(Equal(1))
				Expect(fakePersitentVolumeClient.CreateArgsForCall(0)).To(Equal(
					&v1.PersistentVolume{
						ObjectMeta: metav1.ObjectMeta{
							Name: serviceInstanceKey,
						},
						Spec: v1.PersistentVolumeSpec{
							AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
							Capacity:    v1.ResourceList{v1.ResourceStorage: resource.MustParse("100M")},
							PersistentVolumeSource: v1.PersistentVolumeSource{
								CSI: &v1.CSIPersistentVolumeSource{
									Driver:           "org.cloudfoundry.smb",
									VolumeHandle:     "volume-handle",
									VolumeAttributes: map[string]string{},
								},
							},
						},
					},
				))
			})

			It("should create a persistent volume claim", func() {
				Expect(fakePersitentVolumeClaimClient.CreateCallCount()).To(Equal(1))
				storageClass := ""
				Expect(fakePersitentVolumeClaimClient.CreateArgsForCall(0)).To(Equal(
					&v1.PersistentVolumeClaim{
						ObjectMeta: metav1.ObjectMeta{
							Name: serviceInstanceKey,
						},
						Spec: v1.PersistentVolumeClaimSpec{
							StorageClassName: &storageClass,
							VolumeName:       serviceInstanceKey,
							AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{v1.ResourceStorage: resource.MustParse("1M")},
							},
						},
					},
				))
			})

			Context("when unable to create a persistent volume", func() {
				BeforeEach(func() {
					fakePersitentVolumeClient.CreateReturns(nil, errors.New("K8s ERROR"))
				})

				It("should return a meaningful error", func() {
					Expect(recorder.Code).To(Equal(500))
					bytes, err := ioutil.ReadAll(recorder.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(bytes)).To(Equal("{\"description\":\"K8s ERROR\"}\n"))
				})
			})

			Context("when unable to create a persistent volume claim", func() {
				BeforeEach(func() {
					fakePersitentVolumeClaimClient.CreateReturns(nil, errors.New("K8s ERROR"))
				})

				It("should return a meaningful error", func() {
					Expect(recorder.Code).To(Equal(500))
					bytes, err := ioutil.ReadAll(recorder.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(bytes)).To(Equal("{\"description\":\"K8s ERROR\"}\n"))
				})
			})

			Context("when service instance parameters are not provided", func() {
				BeforeEach(func() {
					request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id" }`))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should allow provisioning and store the new service instance", func() {
					Expect(recorder.Code).To(Equal(201))
					Expect(recorder.Body).To(MatchJSON(`{}`))

				})

				Context("when username is not supplied", func() {

					BeforeEach(func() {
						request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": {"password": "foo"}}`))
						Expect(err).NotTo(HaveOccurred())
					})

					It("should respond with an error", func() {
						Expect(recorder.Code).To(Equal(400))
						Expect(recorder.Body).To(MatchJSON(`{ "description": "both username and password must be provided"}`))

					})
				})

				Context("when password is not supplied", func() {

					BeforeEach(func() {
						request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": {"username": "foo"}}`))
						Expect(err).NotTo(HaveOccurred())
					})

					It("should respond with an error", func() {
						Expect(recorder.Code).To(Equal(400))
						Expect(recorder.Body).To(MatchJSON(`{ "description": "both username and password must be provided"}`))

					})
				})
			})

			Context("when a username and password are supplied", func() {

				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": { "username": "foo", "password": "bar" } }`))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should store the username and password in the PVs volume attributes", func() {
					Expect(fakePersitentVolumeClient.CreateCallCount()).To(Equal(1))
					Expect(fakePersitentVolumeClient.CreateArgsForCall(0)).To(Equal(
						&v1.PersistentVolume{
							ObjectMeta: metav1.ObjectMeta{
								Name: serviceInstanceKey,
							},
							Spec: v1.PersistentVolumeSpec{
								AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
								Capacity:    v1.ResourceList{v1.ResourceStorage: resource.MustParse("100M")},
								PersistentVolumeSource: v1.PersistentVolumeSource{
									CSI: &v1.CSIPersistentVolumeSource{
										Driver:       "org.cloudfoundry.smb",
										VolumeHandle: "volume-handle",
										VolumeAttributes: map[string]string{
											"username": "foo",
											"password": "bar",
										},
									},
								},
							},
						},
					))
				})
			})

			Context("when an invalid username is supplied", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": { "username": 123, "password": "321" } }`))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should respond with an error", func() {
					Expect(recorder.Code).To(Equal(400))
					Expect(recorder.Body).To(MatchJSON(`{ "description": "username must be a string value"}`))
				})
			})

			Context("when an invalid password is supplied", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodPut, "/v2/service_instances/"+serviceInstanceKey, strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "parameters": { "username": "123", "password": 321 } }`))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should respond with an error", func() {
					Expect(recorder.Code).To(Equal(400))
					Expect(recorder.Body).To(MatchJSON(`{ "description": "password must be a string value"}`))
				})
			})
		})

		Describe("#Deprovision endpoint", func() {
			var serviceInstanceKey string
			BeforeEach(func() {
				serviceInstanceKey = randomString(source)

				var err error
				request, err = http.NewRequest(http.MethodDelete, "/v2/service_instances/"+serviceInstanceKey+"?service_id=123&plan_id=plan-id", nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should delete a persistent volume", func() {
				Expect(fakePersitentVolumeClient.DeleteCallCount()).To(Equal(1))
				name, options := fakePersitentVolumeClient.DeleteArgsForCall(0)
				Expect(name).To(Equal(serviceInstanceKey))
				Expect(options).To(Equal(&metav1.DeleteOptions{}))
			})

			It("should delete a persistent volume claim", func() {
				Expect(fakePersitentVolumeClaimClient.DeleteCallCount()).To(Equal(1))
				name, options := fakePersitentVolumeClaimClient.DeleteArgsForCall(0)
				Expect(name).To(Equal(serviceInstanceKey))
				Expect(options).To(Equal(&metav1.DeleteOptions{}))
			})

			Context("when unable to delete a persistent volume", func() {
				BeforeEach(func() {
					fakePersitentVolumeClient.DeleteReturns(errors.New("K8s ERROR"))
				})

				It("should return a meaningful error", func() {
					Expect(recorder.Code).To(Equal(500))
					bytes, err := ioutil.ReadAll(recorder.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(bytes)).To(Equal("{\"description\":\"K8s ERROR\"}\n"))
				})
			})

			Context("when unable to delete a persistent volume claim", func() {
				BeforeEach(func() {
					fakePersitentVolumeClaimClient.DeleteReturns(errors.New("K8s ERROR"))
				})

				It("should return a meaningful error", func() {
					Expect(recorder.Code).To(Equal(500))
					bytes, err := ioutil.ReadAll(recorder.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(bytes)).To(Equal("{\"description\":\"K8s ERROR\"}\n"))
				})
			})
		})

		Describe("#GetInstance endpoint", func() {
			var (
				err                                                               error
				instanceID, val1, val2, val3, key1, key2, key3, serviceID, planID string
			)

			BeforeEach(func() {
				instanceID = randomString(source)
				request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("/v2/service_instances/%s", instanceID), nil)
				Expect(err).NotTo(HaveOccurred())
				request.Header.Add("X-Broker-API-Version", "2.14")
			})

			BeforeEach(func() {
				key1 = "share"
				key2 = "username"
				key3 = "password"
				val1 = randomString(source)
				val2 = randomString(source)
				val3 = randomString(source)
				serviceID = "123"
				planID = "plan-id"

				fakePersitentVolumeClient.GetReturns(&v1.PersistentVolume{
					Spec: v1.PersistentVolumeSpec{
						PersistentVolumeSource: v1.PersistentVolumeSource{
							CSI: &v1.CSIPersistentVolumeSource{
								VolumeAttributes: map[string]string{key1: val1, key2: val2, key3: val3},
							},
						},
					},
				}, nil)
			})

			It("should retrieve a service instance that was earlier provisioned", func() {
				Expect(fakePersitentVolumeClient.GetCallCount()).To(Equal(1))

				instanceIDArg, getOpts := fakePersitentVolumeClient.GetArgsForCall(0)
				Expect(instanceIDArg).To(Equal(instanceID))
				Expect(getOpts).To(Equal(metav1.GetOptions{}))

				Expect(recorder.Code).To(Equal(200))

			})

			It("shows share and username but not password", func() {
				Expect(recorder.Body).To(MatchJSON(
					fmt.Sprintf(`{ "service_id": "%s", "plan_id": "%s", "parameters": { "%s": "%s", "%s": "%s" } }`, serviceID, planID, key1, val1, key2, val2)),
				)
			})

			Context("when no PV exists", func() {
				BeforeEach(func() {
					fakePersitentVolumeClient.GetReturns(nil, errors.New("pv not found"))
				})

				It("Should return an FailureError with a 404 status code", func() {
					Expect(recorder.Code).To(Equal(404))
					Expect(recorder.Body).To(MatchJSON(`{"description": "unable to find service instance"}`))
				})

			})
		})

		Describe("#Bind endpoint", func() {
			var instanceID, bindingID string

			BeforeEach(func() {
				fakePersitentVolumeClient.GetReturns(&v1.PersistentVolume{}, nil)

				instanceID = randomString(source)
				bindingID = randomString(source)
				request, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s", instanceID, bindingID),
					strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "bind_resource": {"app_guid": "456"} }`))
			})

			It("fetches the PV from k8s", func(){
				Expect(fakePersitentVolumeClient.GetCallCount()).To(Equal(1))
				instanceIDArg, optionsArg := fakePersitentVolumeClient.GetArgsForCall(0)
				Expect(instanceIDArg).To(Equal(instanceID))
				Expect(optionsArg).To(Equal(metav1.GetOptions{}))
			})

			Context("given a service instance", func() {
				It("returns a bind response", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(recorder.Code).To(Equal(201))
					Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{"credentials": {}, "volume_mounts": [{"driver": "smb", "container_dir": "/home/vcap/data/", "mode": "rw", "device_type": "shared", "device": {"volume_id": "%s", "mount_config": {"name": "%s"}} }]}`, bindingID, instanceID)))
				})

				Context("given container-dir bind option", func() {
					var mountBindConfig string
					BeforeEach(func() {
						instanceID = randomString(source)
						bindingID = randomString(source)
						mountBindConfig = "/foo/bar"
						request, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s", instanceID, bindingID),
							strings.NewReader(fmt.Sprintf(`{ "service_id": "123", "plan_id": "plan-id", "bind_resource": {"app_guid": "456"}, "parameters": {"mount": "%s"} }`, mountBindConfig)))
					})

					It("should honor that bind option", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(recorder.Code).To(Equal(201))
						Expect(recorder.Body).To(
							MatchJSON(
								fmt.Sprintf(`{"credentials": {}, "volume_mounts": [{"driver": "smb", "container_dir": "%s", "mode": "rw", "device_type": "shared", "device": {"volume_id": "%s", "mount_config": {"name": "%s"}} }]}`,
									mountBindConfig, bindingID, instanceID),
							),
						)
					})
				})

				Context("given invalid parameters", func() {
					BeforeEach(func() {
						instanceID = randomString(source)
						bindingID = randomString(source)
						request, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s", instanceID, bindingID),
							strings.NewReader(`{ "service_id": "123", "plan_id": "plan-id", "bind_resource": {"app_guid": "456"}, "parameters": {"mount": 123} }`))
					})

					It("should return a 422", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(recorder.Code).To(Equal(422))
						Expect(recorder.Body).To(
							MatchJSON(
								`{"description": "The format of the parameters is not valid JSON"}`,
							),
						)
					})
				})
			})

			Context("given the service instance doesnt exist", func() {
				BeforeEach(func() {
					fakePersitentVolumeClient.GetReturns(&v1.PersistentVolume{}, errors.New("pv not found"))
				})

				It("should return an error", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(recorder.Code).To(Equal(404))
				})
			})
		})

		Describe("Unbind endpoint", func() {
			var instanceID, bindingID string

			BeforeEach(func() {
				fakePersitentVolumeClient.GetReturns(&v1.PersistentVolume{}, nil)
				instanceID = randomString(source)
				bindingID = randomString(source)
				request, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s?service_id=123&plan_id=plan-id", instanceID, bindingID), nil)
			})

			It("returns 200", func() {
				Expect(fakePersitentVolumeClient.GetCallCount()).To(Equal(1))
				instanceIDArg, getOpts := fakePersitentVolumeClient.GetArgsForCall(0)
				Expect(instanceIDArg).To(Equal(instanceID))
				Expect(getOpts).To(Equal(metav1.GetOptions{}))

				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Code).To(Equal(200))
			})

			Context("given the service instance doesnt exist", func() {
				BeforeEach(func() {
					fakePersitentVolumeClient.GetReturns(&v1.PersistentVolume{}, errors.New("pv does not exist"))
				})

				It("should return an error", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(recorder.Code).To(Equal(410))
				})
			})

		})
	})
})

func randomString(sourceSeededByGinkgo rand.Source) string {
	return strconv.Itoa(rand.New(sourceSeededByGinkgo).Int())
}
