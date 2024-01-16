package e2e

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	CI "github.com/terraform-redhat/terraform-provider-rhcs/tests/ci"
	CON "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
	. "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/log"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/openshift"
)

var _ = Describe("RHCS Provider Test", func() {
	Describe("Create cluster test", func() {
		It("CreateClusterByProfile", CI.Day1Prepare,
			func() {

				// Generate/build cluster by profile selected
				profile := CI.LoadProfileYamlFileByENV()
				clusterID, err := CI.CreateRHCSClusterByProfile(token, profile)
				Expect(err).ToNot(HaveOccurred())
				Expect(clusterID).ToNot(BeEmpty())
				//TODO: implement waiter for  the private cluster once bastion is implemented
				if CON.GetEnvWithDefault(CON.WaitOperators, "false") == "true" && !profile.Private {
					// WaitClusterOperatorsToReadyStatus will wait for cluster operators ready
					timeout := 60
					timeoutMin := time.Duration(timeout)
					console, err := openshift.NewConsole(clusterID, CI.RHCSConnection)
					if err != nil {
						Logger.Warnf("Got error %s when config the openshift console. Return without waiting for operators ready", err.Error())
						return
					}
					_, err = openshift.RetryCMDRun(fmt.Sprintf("oc wait clusteroperators --all --for=condition=Progressing=false --kubeconfig %s --timeout %dm", console.KubePath, timeout), timeoutMin)
					Expect(err).ToNot(HaveOccurred())
				}
			})
	})
})
