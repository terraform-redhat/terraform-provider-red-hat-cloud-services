package e2e

import (

	// nolint

	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ci "github.com/terraform-redhat/terraform-provider-rhcs/tests/ci"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/cms"
	con "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
	exe "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/exec"
)

var _ = Describe("Create cluster autoscaler", func() {
	defer GinkgoRecover()

	var caService *exe.ClusterAutoscalerService
	var clusterAutoScalerBodyForRecreate *cmv1.ClusterAutoscaler
	var clusterAutoscalerStatusBefore int

	BeforeEach(func() {
		caService = exe.NewClusterAutoscalerService(con.ClusterAutoscalerDir)
		caRetrieveBody, _ := cms.RetrieveClusterAutoscaler(ci.RHCSConnection, clusterID)
		clusterAutoscalerStatusBefore = caRetrieveBody.Status()
		if clusterAutoscalerStatusBefore == http.StatusOK {
			clusterAutoScalerBodyForRecreate = caRetrieveBody.Body()
		}
	})
	AfterEach(func() {
		By("Recover clusterautoscaler")
		clusterAutoscalerAfter, _ := cms.RetrieveClusterAutoscaler(ci.RHCSConnection, clusterID)
		if (clusterAutoscalerAfter.Status() == clusterAutoscalerStatusBefore) && clusterAutoscalerStatusBefore != http.StatusNotFound {
			recreateAutoscaler, err := cms.PatchClusterAutoscaler(ci.RHCSConnection, clusterID, clusterAutoScalerBodyForRecreate)
			Expect(err).NotTo(HaveOccurred())
			Expect(recreateAutoscaler.Status()).To(Equal(http.StatusOK))
		} else if clusterAutoscalerAfter.Status() == http.StatusOK && clusterAutoscalerStatusBefore == http.StatusNotFound {
			deleteAutoscaler, err := cms.DeleteClusterAutoscaler(ci.RHCSConnection, clusterID)
			Expect(err).NotTo(HaveOccurred())
			Expect(deleteAutoscaler.Status()).To(Equal(http.StatusNoContent))
		} else if clusterAutoscalerAfter.Status() == http.StatusNotFound && clusterAutoscalerStatusBefore == http.StatusOK {
			recreateAutoscaler, err := cms.CreateClusterAutoscaler(ci.RHCSConnection, clusterID, clusterAutoScalerBodyForRecreate)
			Expect(err).NotTo(HaveOccurred())
			Expect(recreateAutoscaler.Status()).To(Equal(http.StatusCreated))
		}
	})
	It("works with ocm terraform provider - [id:69137]", ci.Day2, ci.High, ci.NonHCPCluster, ci.FeatureClusterautoscaler, func() {
		By("Delete clusterautoscaler when it exists in cluster")
		if clusterAutoscalerStatusBefore == http.StatusOK {
			caDeleteBody, err := cms.DeleteClusterAutoscaler(ci.RHCSConnection, clusterID)
			Expect(err).NotTo(HaveOccurred())
			Expect(caDeleteBody.Status()).To(Equal(http.StatusNoContent))
		}

		By("Create clusterautoscaler")
		max := 1
		min := 0
		resourceRange := &exe.ResourceRange{
			Max: max,
			Min: min,
		}
		maxNodesTotal := 10
		resourceLimits := &exe.ResourceLimits{
			Cores:         resourceRange,
			MaxNodesTotal: maxNodesTotal,
			Memory:        resourceRange,
		}
		delayAfterAdd := "3h"
		delayAfterDelete := "3h"
		delayAfterFailure := "3h"
		unneededTime := "1h"
		utilizationThreshold := "0.5"
		enabled := true
		scaleDown := &exe.ScaleDown{
			DelayAfterAdd:        delayAfterAdd,
			DelayAfterDelete:     delayAfterDelete,
			DelayAfterFailure:    delayAfterFailure,
			UnneededTime:         unneededTime,
			UtilizationThreshold: utilizationThreshold,
			Enabled:              enabled,
		}
		balanceSimilarNodeGroups := true
		skipNodesWithLocalStorage := true
		logVerbosity := 1
		maxPodGracePeriod := 10
		podPriorityThreshold := -10
		ignoreDaemonsetsUtilization := true
		maxNodeProvisionTime := "1h"
		balancingIgnoredLabels := []string{"l1", "l2"}
		ClusterAutoscalerArgs := &exe.ClusterAutoscalerArgs{
			Cluster:                     clusterID,
			BalanceSimilarNodeGroups:    balanceSimilarNodeGroups,
			SkipNodesWithLocalStorage:   skipNodesWithLocalStorage,
			LogVerbosity:                logVerbosity,
			MaxPodGracePeriod:           maxPodGracePeriod,
			PodPriorityThreshold:        podPriorityThreshold,
			IgnoreDaemonsetsUtilization: ignoreDaemonsetsUtilization,
			MaxNodeProvisionTime:        maxNodeProvisionTime,
			BalancingIgnoredLabels:      balancingIgnoredLabels,
			ResourceLimits:              resourceLimits,
			ScaleDown:                   scaleDown,
		}
		_, err = caService.Apply(ClusterAutoscalerArgs, false)
		Expect(err).ToNot(HaveOccurred())
		_, err = caService.Output()
		Expect(err).ToNot(HaveOccurred())

		By("Verify the parameters of the createdautoscaler")
		caOut, err := caService.Output()
		Expect(err).ToNot(HaveOccurred())
		caResponseBody, err := cms.RetrieveClusterAutoscaler(ci.RHCSConnection, clusterID)
		Expect(err).ToNot(HaveOccurred())
		Expect(caResponseBody.Body().BalanceSimilarNodeGroups()).To(Equal(caOut.BalanceSimilarNodeGroups))
		Expect(caResponseBody.Body().SkipNodesWithLocalStorage()).To(Equal(caOut.SkipNodesWithLocalStorage))
		Expect(caResponseBody.Body().LogVerbosity()).To(Equal(caOut.LogVerbosity))
		Expect(caResponseBody.Body().MaxPodGracePeriod()).To(Equal(caOut.MaxPodGracePeriod))
		Expect(caResponseBody.Body().PodPriorityThreshold()).To(Equal(caOut.PodPriorityThreshold))
		Expect(caResponseBody.Body().IgnoreDaemonsetsUtilization()).To(Equal(caOut.IgnoreDaemonsetsUtilization))
		Expect(caResponseBody.Body().MaxNodeProvisionTime()).To(Equal(caOut.MaxNodeProvisionTime))
		Expect(caResponseBody.Body().BalancingIgnoredLabels()).To(Equal(caOut.BalancingIgnoredLabels))
		Expect(caResponseBody.Body().ResourceLimits().MaxNodesTotal()).To(Equal(caOut.MaxNodesTotal))
		Expect(caResponseBody.Body().ScaleDown().DelayAfterAdd()).To(Equal(caOut.DelayAfterAdd))
		Expect(caResponseBody.Body().ScaleDown().DelayAfterDelete()).To(Equal(caOut.DelayAfterDelete))
		Expect(caResponseBody.Body().ScaleDown().DelayAfterFailure()).To(Equal(caOut.DelayAfterFailure))
		Expect(caResponseBody.Body().ScaleDown().UnneededTime()).To(Equal(caOut.UnneededTime))
		Expect(caResponseBody.Body().ScaleDown().UtilizationThreshold()).To(Equal(caOut.UtilizationThreshold))
		Expect(caResponseBody.Body().ScaleDown().Enabled()).To(Equal(caOut.Enabled))
	})
})
