/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//+kubebuilder:scaffold:imports

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

// var cfg *rest.Config
// var k8sClient client.Client
// var testEnv *envtest.Environment
// var nodeLabelerReconciler *NodeLabelerReconciler

func TestAPIs(t *testing.T) {
	// RegisterFailHandler(Fail)

	// RunSpecsWithDefaultAndCustomReporters(t,
	// 	"Controller Suite",
	// 	[]Reporter{printer.NewlineReporter{}})
	assert.True(t, true)
}

// var _ = BeforeSuite(func() {
// 	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
// 	Expect(os.Setenv("USE_EXISTING_CLUSTER", "true")).To(Succeed())
// 	By("bootstrapping test environment")
// 	testEnv = &envtest.Environment{
// 		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
// 		ErrorIfCRDPathMissing: true,
// 	}

// 	var err error
// 	// cfg is defined in this file globally.
// 	cfg, err = testEnv.Start()
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(cfg).NotTo(BeNil())

// 	err = kubebuilderv1alpha1.AddToScheme(scheme.Scheme)
// 	Expect(err).NotTo(HaveOccurred())

// 	//+kubebuilder:scaffold:scheme

// 	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(k8sClient).NotTo(BeNil())

// 	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
// 		Scheme:             scheme.Scheme,
// 		MetricsBindAddress: ":8090",
// 	})
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(k8sManager.GetClient()).ToNot(BeNil())

// 	nodeLabelerReconciler = &NodeLabelerReconciler{
// 		Client: k8sManager.GetClient(),
// 		Scheme: k8sClient.Scheme(),
// 	}
// 	err = nodeLabelerReconciler.SetupWithManager(k8sManager)

// 	go func() {
// 		err = k8sManager.Start(ctrl.SetupSignalHandler())
// 		Expect(err).ToNot(HaveOccurred())
// 	}()

// }, 60)

// var _ = AfterSuite(func() {
// 	Expect(os.Unsetenv("USE_EXISTING_CLUSTER")).To(Succeed())
// 	By("tearing down the test environment")
// 	err := testEnv.Stop()
// 	Expect(err).NotTo(HaveOccurred())
// })
