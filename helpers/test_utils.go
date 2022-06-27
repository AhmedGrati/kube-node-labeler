package helpers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
func getBasicNode() *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "basic-node",
			Labels: map[string]string{
				"os":              "linux",
				"number-of-years": "2",
				"ip-address":      "127.0.0.1",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: []corev1.Taint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: corev1.TaintEffectNoExecute,
				},
			},
		},
	}
}
func generateNodeListForRegexTest() *corev1.NodeList {
	return &corev1.NodeList{
		Items: []corev1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "minikube1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "2minikube2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "3minikube3",
				},
			},
		},
	}
}
