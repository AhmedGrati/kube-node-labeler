package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)
var _ = Describe("NodeLabeler Controller Test", func() {
	Context("Simple NodeLabeler", func() {
		nodeLabelerObject := generateSampleNodeLabelerObject()
		It("Should create a node labeler object successfully", func () {
			Expect(k8sClient.Create(context.Background(), nodeLabelerObject)).Should(Succeed())
		})
	})
})