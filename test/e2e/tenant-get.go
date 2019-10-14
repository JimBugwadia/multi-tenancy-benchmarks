package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	gomega "github.com/onsi/gomega"

	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/test/e2e/framework"
	ginkgowrapper "k8s.io/kubernetes/test/e2e/framework/ginkgowrapper"
)

const (
	expectedVal      = "Error from server (Forbidden)"
	adminkubeconfig  = "manifest/admin-kubeconfig"
	tenantkubeconfig = "manifest/tenant-kubeconfig"
)

var _ = framework.KubeDescribe("test tenant permission", func() {
	var resourceList string
	var err error

	framework.KubeDescribe("test tenant get none namespaced resource", func() {
		ginkgo.BeforeEach(func() {
			ginkgo.By("get none namespaced api-resources")
			os.Setenv("KUBECONFIG", adminkubeconfig)
			nsdFlag := fmt.Sprintf("--namespaced=false")
			outputFlag := fmt.Sprintf("-o=name")

			resourceList, err = framework.RunKubectl("api-resources", nsdFlag, outputFlag)
			framework.ExpectNoError(err)
		})

		// ginkgo.JustBeforeEach(func() {
		// 	os.Setenv("KUBECONFIG", tenantkubeconfig)
		// })

		framework.KubeDescribe("tenant admin", func() {
			// mkpath := func(file string) string {
			// 	return filepath.Join(manifestPath, file)
			// }

			ginkgo.BeforeEach(func() {
				os.Setenv("KUBECONFIG", tenantkubeconfig)
			})

			ginkgo.It("tenant admin cannot get none namespaced resources ", func() {

				resources := strings.Fields(resourceList)
				for _, resource := range resources {
					_, err1 := framework.LookForString(expectedVal, time.Minute, func() string {
						_, err := framework.RunKubectl("get", resource)
						return err.Error()
					})

					framework.ExpectNoError(err1)
				}
			})
		})
	})
})

func RunE2ETests(t *testing.T) {
	logs.InitLogs()
	defer logs.FlushLogs()

	gomega.RegisterFailHandler(ginkgowrapper.Fail)
	ginkgo.RunSpecs(t, "Multi-Tenancy Benchmarks")
}