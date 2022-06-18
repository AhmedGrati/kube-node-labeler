package controllers

import (
	"kube-node-labeler/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)
const(
	NodeLabelerName = "simple-node-labeler"
	NodeLabelerNamespace = "default"
	NodeLabelerKind = "NodeLabeler"
	NodeLabelerAPIVersion = "alpha1v1"
)
var NodeLabelerLabels = map[string]string{
		"type": "unit-test",
	}
var NodeLabelerAnnotations = map[string]string {
	"annotation1": "value1",
}
func generateSampleTypeMeta() *metav1.TypeMeta {
	return &metav1.TypeMeta{
		Kind: NodeLabelerKind,
		APIVersion: NodeLabelerAPIVersion,
	}
}

func generateSampleObjectMeta() *metav1.ObjectMeta {
	return &metav1.ObjectMeta{
		Name: NodeLabelerName,
		Namespace: NodeLabelerNamespace,
		Labels: NodeLabelerLabels,
		Annotations: NodeLabelerAnnotations,
	}
}

func generateMatchExpressions() *[]corev1.NodeSelectorRequirement {
	return &[]corev1.NodeSelectorRequirement{
		{
			Key: "beta.kubernetes.io/arch",
			Operator: corev1.NodeSelectorOperator("In"),
			Values: []string{
				"amd64", "arch",
			},
		},
	}
}

func generateSampleNodeLabelerSpec() *v1alpha1.NodeLabelerSpec {
	return &v1alpha1.NodeLabelerSpec{
		NodeSelector: corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: *generateMatchExpressions(),
				},
			},
		},
	}
}

func generateSampleNodeLabelerObject() *v1alpha1.NodeLabeler {
	return &v1alpha1.NodeLabeler{
		TypeMeta: *generateSampleTypeMeta(),
		ObjectMeta: *generateSampleObjectMeta(),
		Spec: *generateSampleNodeLabelerSpec(),
	}
}