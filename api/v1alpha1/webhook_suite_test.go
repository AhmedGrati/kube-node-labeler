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

package v1alpha1

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

// func TestAPIs(t *testing.T) {
// 	RegisterFailHandler(Fail)

// 	RunSpecsWithDefaultAndCustomReporters(t,
// 		"Webhook Suite",
// 		[]Reporter{printer.NewlineReporter{}})
// }

// var _ = BeforeSuite(func() {
// 	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

// 	ctx, cancel = context.WithCancel(context.TODO())

// 	By("bootstrapping test environment")
// 	testEnv = &envtest.Environment{
// 		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
// 		ErrorIfCRDPathMissing: false,
// 		WebhookInstallOptions: envtest.WebhookInstallOptions{
// 			Paths: []string{filepath.Join("..", "..", "config", "webhook")},
// 		},
// 	}

// 	var err error
// 	// cfg is defined in this file globally.
// 	cfg, err = testEnv.Start()
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(cfg).NotTo(BeNil())

// 	scheme := runtime.NewScheme()
// 	err = AddToScheme(scheme)
// 	Expect(err).NotTo(HaveOccurred())

// 	err = admissionv1beta1.AddToScheme(scheme)
// 	Expect(err).NotTo(HaveOccurred())

// 	//+kubebuilder:scaffold:scheme

// 	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(k8sClient).NotTo(BeNil())

// 	// start webhook server using Manager
// 	webhookInstallOptions := &testEnv.WebhookInstallOptions
// 	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
// 		Scheme:             scheme,
// 		Host:               webhookInstallOptions.LocalServingHost,
// 		Port:               webhookInstallOptions.LocalServingPort,
// 		CertDir:            webhookInstallOptions.LocalServingCertDir,
// 		LeaderElection:     false,
// 		MetricsBindAddress: "0",
// 	})
// 	Expect(err).NotTo(HaveOccurred())

// 	err = (&NodeLabeler{}).SetupWebhookWithManager(mgr)
// 	Expect(err).NotTo(HaveOccurred())

// 	//+kubebuilder:scaffold:webhook

// 	go func() {
// 		defer GinkgoRecover()
// 		err = mgr.Start(ctx)
// 		Expect(err).NotTo(HaveOccurred())
// 	}()

// 	// wait for the webhook server to get ready
// 	dialer := &net.Dialer{Timeout: time.Second}
// 	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
// 	Eventually(func() error {
// 		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
// 		if err != nil {
// 			return err
// 		}
// 		conn.Close()
// 		return nil
// 	}).Should(Succeed())

// }, 60)

// var _ = AfterSuite(func() {
// 	cancel()
// 	By("tearing down the test environment")
// 	err := testEnv.Stop()
// 	Expect(err).NotTo(HaveOccurred())
// })

func getNodeLabelerWithoutSize() *NodeLabeler {
	return &NodeLabeler{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-labeler",
		},
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getNodeLabelerWithInvalidName() *NodeLabeler {
	return &NodeLabeler{
		ObjectMeta: metav1.ObjectMeta{
			Name: randStringBytes(validationutils.DNS1035LabelMaxLength - 10),
		},
	}
}

func TestDefaultingWebhook(t *testing.T) {
	r := getNodeLabelerWithoutSize()
	r.Default()
	assert.Equal(t, *r.Spec.Size, int(0))
	assert.NotEqual(t, r.Spec.Size, nil)
	*r.Spec.Size = 2
	r.Default()
	assert.Equal(t, *r.Spec.Size, int(2))
}

/*
After Triggering error in each field we correct the error before moving to test the next field
*/
func TestCreateValidationWebhook(t *testing.T) {
	nodeLabeler := getNodeLabelerWithInvalidName()
	err := nodeLabeler.validateNodeLabelerName()
	assert.Equal(t, err.Field, "metadata.name")

	nodeLabeler.ObjectMeta.Name = randStringBytes(10)
	nodeLabeler.Default()

	// Invalid Size
	*nodeLabeler.Spec.Size = -22
	err = nodeLabeler.validateNodeLabelerSpec()
	assert.Equal(t, err.Field, "spec.size")

	*nodeLabeler.Spec.Size = 3
	// Invalid Regex
	nodeLabeler.Spec.NodeNamePatterns = []string{`^\/(?!\/)(.*?)`}
	err = nodeLabeler.validateNodeLabelerSpec()
	assert.Equal(t, err.Field, "spec.nodeNamePatterns")
}
