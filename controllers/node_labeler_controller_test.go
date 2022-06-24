package controllers

import (
	"context"
	"kube-node-labeler/api/v1alpha1"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestNodeLabelerController(t *testing.T) {
	name := "custom-node-labeler"
	nodelabeler := generateSampleNodeLabelerObject()
	objs := []runtime.Object{nodelabeler}
	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeBuilder.GroupVersion, nodelabeler)
	cl := fake.NewClientBuilder().WithRuntimeObjects(objs...).WithScheme(s).Build()
	r := &NodeLabelerReconciler{Client: cl, Scheme: s}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: "",
		},
	}
	res, err := r.Reconcile(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.Requeue, "We should not requeue")
}
