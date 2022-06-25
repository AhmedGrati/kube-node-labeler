package controllers

import (
	"context"
	"kube-node-labeler/api/v1alpha1"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func getNodeLabelerReconciler(objs []runtime.Object) (*NodeLabelerReconciler, client.Client) {
	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.GroupVersion, &v1alpha1.NodeLabeler{})
	cl := fake.NewClientBuilder().WithScheme(s).WithRuntimeObjects(objs...).Build()
	r := New(cl, s)
	return r, cl
}
