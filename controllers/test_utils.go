package controllers

import (
	"kube-node-labeler/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NodeLabelerName       = "simple-node-labeler"
	NodeLabelerNamespace  = "default"
	NodeLabelerKind       = "NodeLabeler"
	NodeLabelerAPIVersion = "alpha1v1"
)

var NodeLabelerLabels = map[string]string{
	"type": "unit-test",
}
var NodeLabelerAnnotations = map[string]string{
	"annotation1": "value1",
}

var LabelsToMerge = map[string]string{
	"merge-label": "true",
	"test-label":  "true",
}

var AnnotationsToMerge = map[string]string{
	"merge-annotation": "true",
}

var TaintsToMerge = []corev1.Taint{
	{
		Key:    "key1",
		Value:  "value1",
		Effect: corev1.TaintEffectNoSchedule,
	},
	{
		Key:    "key2",
		Value:  "value2",
		Effect: corev1.TaintEffectNoExecute,
	},
}

var LabelsToOverwrite = map[string]string{
	"overwrite-label": "true",
	"merge-label":     "false",
}

var AnnotationsToOverwrite = map[string]string{
	"overwrite-annotation": "true",
	"merge-annotation":     "false",
}

var TaintsToOverwrite = []corev1.Taint{
	{
		Key:    "key1",
		Value:  "value1",
		Effect: corev1.TaintEffectPreferNoSchedule,
	},
}

func generateSampleTypeMeta() *metav1.TypeMeta {
	return &metav1.TypeMeta{
		Kind:       NodeLabelerKind,
		APIVersion: NodeLabelerAPIVersion,
	}
}

func generateSampleObjectMeta() *metav1.ObjectMeta {
	return &metav1.ObjectMeta{
		Name:        NodeLabelerName,
		Namespace:   NodeLabelerNamespace,
		Labels:      NodeLabelerLabels,
		Annotations: NodeLabelerAnnotations,
	}
}

func generateMatchExpressions() *[]corev1.NodeSelectorRequirement {
	return &[]corev1.NodeSelectorRequirement{
		{
			Key:      "beta.kubernetes.io/arch",
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
		Merge:     *generateSampleMergeSpec(),
		Overwrite: *generateSampleOverwriteSpec(),
	}
}

func generateSampleOverwriteSpec() *v1alpha1.OverwriteSpec {
	return &v1alpha1.OverwriteSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      LabelsToOverwrite,
			Annotations: AnnotationsToOverwrite,
		},
		NodeSpec: corev1.NodeSpec{
			Taints: TaintsToOverwrite,
		},
	}
}

func generateSampleMergeSpec() *v1alpha1.MergeSpec {
	return &v1alpha1.MergeSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      LabelsToMerge,
			Annotations: AnnotationsToMerge,
		},
		NodeSpec: corev1.NodeSpec{
			Taints: TaintsToMerge,
		},
	}
}

func generateSampleNodeLabelerObject() *v1alpha1.NodeLabeler {
	return &v1alpha1.NodeLabeler{
		TypeMeta:   *generateSampleTypeMeta(),
		ObjectMeta: *generateSampleObjectMeta(),
		Spec:       *generateSampleNodeLabelerSpec(),
	}
}

// func generateWrongNodeLabelerObjectMeta() *metav1.ObjectMeta {
// 	return &metav1.ObjectMeta{
// 		Name:      NodeLabelerName,
// 		Namespace: "custom-namespace",
// 	}
// }

// func generateWrongNodeLabelerObject() *v1alpha1.NodeLabeler {
// 	return &v1alpha1.NodeLabeler{
// 		TypeMeta:   *generateSampleTypeMeta(),
// 		ObjectMeta: *generateWrongNodeLabelerObjectMeta(),
// 		Spec:       *generateSampleNodeLabelerSpec(),
// 	}
// }

func getNode() *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
			Labels: map[string]string{
				"os":                      "linux",
				"number-of-years":         "2",
				"ip-address":              "127.0.0.1",
				"beta.kubernetes.io/arch": "arch",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: []corev1.Taint{
				{
					Key:    "randomkey",
					Value:  "randomvalue",
					Effect: corev1.TaintEffectNoExecute,
				},
			},
		},
	}
}
