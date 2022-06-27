package controllers

import (
	"context"
	"fmt"
	"kube-node-labeler/api/v1alpha1"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestNodeLabelerSuccessfullCreation(t *testing.T) {
	nodeLabeler := generateSampleNodeLabelerObject()
	objs := []runtime.Object{nodeLabeler}

	r, cl := getNodeLabelerReconciler(objs)
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      nodeLabeler.Name,
			Namespace: nodeLabeler.Namespace,
		},
	}
	res, err := r.Reconcile(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.Requeue, "We should not requeue")

	nl := &v1alpha1.NodeLabeler{}
	err = cl.Get(context.Background(), types.NamespacedName{
		Name:      nodeLabeler.Name,
		Namespace: nodeLabeler.Namespace,
	},
		nl,
	)

	assert.Equal(t, nl.Name, nodeLabeler.Name)
	assert.NoError(t, err)
}

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

func TestNodesManagement(t *testing.T) {
	nodeLabeler := generateSampleNodeLabelerObject()
	objs := []runtime.Object{nodeLabeler}
	r, _ := getNodeLabelerReconciler(objs)
	node := *getNode()
	nodes := &corev1.NodeList{
		Items: []corev1.Node{
			node,
		},
	}
	nodeLabelerSpec := generateSampleNodeLabelerSpec()
	managedNodes, err := r.ManageNodes(nodes, *nodeLabelerSpec, len(nodes.Items))
	assert.NoError(t, err)
	updatedNode := managedNodes.Items[0]
	// verify that managed node contains the desired labels
	assert.Equal(t, len(updatedNode.Labels), len(node.Labels)+3)
	for k := range LabelsToMerge {
		assert.Contains(t, updatedNode.Labels, k)
	}
	for k := range LabelsToOverwrite {
		assert.Contains(t, updatedNode.Labels, k)
	}
	assert.Equal(t, updatedNode.Labels["merge-label"], "false")
	assert.Equal(t, updatedNode.Labels["overwrite-label"], "true")
	assert.Equal(t, updatedNode.Labels["test-label"], "true")

	// verify that managed node contains the desired annotations
	assert.Equal(t, len(updatedNode.Annotations), len(node.Annotations)+2)
	for k := range AnnotationsToMerge {
		assert.Contains(t, updatedNode.Annotations, k)
	}
	for k := range AnnotationsToOverwrite {
		assert.Contains(t, updatedNode.Annotations, k)
	}
	assert.Equal(t, updatedNode.Annotations["merge-annotation"], "false")
	assert.Equal(t, updatedNode.Annotations["overwrite-annotation"], "true")

	// verify that managed node contains the desired taints
	assert.Equal(t, len(updatedNode.Spec.Taints), len(node.Spec.Taints)+2)
	fmt.Println(updatedNode.Spec.Taints)
	desiredTaints := append(node.Spec.Taints, []corev1.Taint{
		{
			Key:    "key1",
			Value:  "value1",
			Effect: corev1.TaintEffectPreferNoSchedule,
		},
		{
			Key:    "key2",
			Value:  "value2",
			Effect: corev1.TaintEffectNoExecute,
		},
	}...,
	)
	assert.Equal(t, updatedNode.Spec.Taints, desiredTaints)
}

func getNodeLabelerReconciler(objs []runtime.Object) (*NodeLabelerReconciler, client.Client) {
	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.GroupVersion, &v1alpha1.NodeLabeler{})
	cl := fake.NewClientBuilder().WithScheme(s).WithRuntimeObjects(objs...).Build()
	r := New(cl, s)
	return r, cl
}

func TestRegisterWithManager(t *testing.T) {
	t.Skip("this test requires a real cluster, otherwise the GetConfigOrDie will die")

	// prepare
	mgr, err := manager.New(k8sconfig.GetConfigOrDie(), manager.Options{})
	require.NoError(t, err)
	nodeLabeler := generateSampleNodeLabelerObject()
	objs := []runtime.Object{nodeLabeler}

	r, _ := getNodeLabelerReconciler(objs)

	// test
	err = r.SetupWithManager(mgr)

	// verify
	assert.NoError(t, err)
}
