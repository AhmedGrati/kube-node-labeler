package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)
var _ = Describe("NodeLabeler Controller Test", func() {
	Context("Simple NodeLabeler", func() {
		nodeLabelerObject := generateSampleNodeLabelerObject()
		wrongNodeLabelerObject := generateWrongNodeLabelerObject()
		It("Should create a node labeler object successfully", func () {
			Expect(k8sClient.Create(context.Background(), nodeLabelerObject)).Should(Succeed())
		})
		It("Should delete a node labeler object successfully", func() {
			Expect(k8sClient.Delete(context.Background(), nodeLabelerObject)).Should(Succeed())
		})
		It("Should Update a node labeler object successfully", func() {
			nodeLabelerObject.ObjectMeta.Name = "another-node-labeler-name"
			Expect(k8sClient.Update(context.Background(), nodeLabelerObject, &client.UpdateOptions{})).Should(Succeed())
		})
		It("Should fail creating a node labeler", func() {
			Expect(k8sClient.Create(context.Background(), wrongNodeLabelerObject)).ShouldNot(Succeed())
		})
	})
})