package e2e

import (

	// nolint

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmsv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/terraform-redhat/terraform-provider-rhcs/tests/ci"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/cms"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/exec"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/helper"
	. "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/log"
)

var internalListeningMethod = "internal"
var externalListeningMethod = "external"

var _ = Describe("HCP Ingress", ci.FeatureIngress, ci.Day2, func() {

	var (
		err            error
		ingressBefore  *cmsv1.Ingress
		ingressService exec.IngressService
		ingressArgs    *exec.IngressArgs
	)

	initializeIngressArgs := func() {
		ingressArgs, err = ingressService.ReadTFVars()
		Expect(err).ToNot(HaveOccurred())
		if ingressArgs.Cluster == nil {
			ingressArgs.Cluster = helper.StringPointer(clusterID)
		}
	}

	BeforeEach(func() {
		profile := ci.LoadProfileYamlFileByENV()
		if !profile.GetClusterType().HCP {
			Skip("Test can run only on Hosted cluster")
		}

		ingressBefore, err = cms.RetrieveClusterIngress(ci.RHCSConnection, clusterID)
		Expect(err).ToNot(HaveOccurred())

		ingressService, err = exec.NewIngressService(constants.HCPIngressDir)
		Expect(err).ToNot(HaveOccurred())

		initializeIngressArgs()
	})

	AfterEach(func() {
		ingressArgs.ListeningMethod = helper.StringPointer(string(ingressBefore.Listening()))
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).ToNot(HaveOccurred())
	})

	It("can be edited - [id:72517]",
		ci.High,
		func() {
			By("Set Listening method to internal")
			ingressArgs.ListeningMethod = helper.StringPointer(internalListeningMethod)
			_, err = ingressService.Apply(ingressArgs)
			Expect(err).ToNot(HaveOccurred())

			By("Verify Cluster Ingress")
			ingress, err := cms.RetrieveClusterIngress(ci.RHCSConnection, clusterID)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(ingress.Listening())).To(Equal("internal"))

			By("Set Listening method to external")
			ingressArgs.ListeningMethod = helper.StringPointer(externalListeningMethod)
			_, err = ingressService.Apply(ingressArgs)
			Expect(err).ToNot(HaveOccurred())

			By("Verify Cluster Ingress")
			ingress, err = cms.RetrieveClusterIngress(ci.RHCSConnection, clusterID)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(ingress.Listening())).To(Equal("external"))

			By("Destroy Cluster Ingress")
			_, err = ingressService.Destroy()
			Expect(err).ToNot(HaveOccurred())

			By("Verify Cluster Ingress is still present")
			ingress, err = cms.RetrieveClusterIngress(ci.RHCSConnection, clusterID)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(ingress.Listening())).To(Equal("external"))
		})

	It("validate edit - [id:72520]", ci.Medium, func() {
		By("Initialize ingress state")
		initializeIngressArgs()
		ingressArgs.ListeningMethod = helper.StringPointer(string(ingressBefore.Listening()))
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).ToNot(HaveOccurred())

		By("Try to edit with empty cluster")
		initializeIngressArgs()
		ingressArgs.Cluster = helper.EmptyStringPointer
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Attribute cluster cluster ID may not be empty/blank string"))

		By("Try to edit cluster with other cluster ID")
		clustersResp, err := cms.ListClusters(ci.RHCSConnection)
		Expect(err).ToNot(HaveOccurred())
		var otherClusterID string
		for _, cluster := range clustersResp.Items().Slice() {
			if cluster.ID() != clusterID && cluster.Status().State() == cmsv1.ClusterStateReady {
				otherClusterID = cluster.ID()
				break
			}
		}
		if otherClusterID != "" {
			initializeIngressArgs()
			ingressArgs.Cluster = helper.StringPointer(otherClusterID)
			_, err = ingressService.Apply(ingressArgs)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Attribute cluster, cannot be changed from"))
		} else {
			Logger.Info("No other cluster accessible for testing this change")
		}

		By("Try to edit cluster field with wrong value")
		initializeIngressArgs()
		ingressArgs.Cluster = helper.StringPointer("wrong")
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Cluster 'wrong' not"))

		By("Try to edit with empty listening_method")
		initializeIngressArgs()
		ingressArgs.ListeningMethod = helper.EmptyStringPointer
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Expected a valid param"))
		Expect(err.Error()).To(ContainSubstring("Options are"))

		By("Try to edit with wrong listening_method")
		initializeIngressArgs()
		ingressArgs.ListeningMethod = helper.StringPointer("wrong")
		_, err = ingressService.Apply(ingressArgs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Expected a valid param"))
		Expect(err.Error()).To(ContainSubstring("Options are"))
	})
})
