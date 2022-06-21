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
	"context"
	"fmt"
	"kube-node-labeler/api/v1alpha1"
	kubebuilderv1alpha1 "kube-node-labeler/api/v1alpha1"
	"kube-node-labeler/helpers"
	"reflect"

	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	MergeStrategy     = "merge"
	OverwriteStrategy = "overwrite"
)

// NodeLabelerReconciler reconciles a NodeLabeler object
type NodeLabelerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubebuilder.kube.node.labeler.io,resources=nodelabelers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubebuilder.kube.node.labeler.io,resources=nodelabelers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubebuilder.kube.node.labeler.io,resources=nodelabelers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NodeLabeler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile

func (r *NodeLabelerReconciler) getAllNodes(ctx context.Context) corev1.NodeList {
	nodes := &corev1.NodeList{}
	opts := []client.ListOption{}
	r.List(ctx, nodes, opts...)
	return *nodes
}

func (r *NodeLabelerReconciler) AssignAttributesToNodes(node *corev1.Node, l metav1.ObjectMeta, spec corev1.NodeSpec, strategyFunc func(*mergo.Config)) {
	cop := node.DeepCopy()
	if err := mergo.Merge(&cop.ObjectMeta, l, strategyFunc); err != nil {
		_ = fmt.Errorf("ERROR: %s", err)
	}
	if err := mergo.Merge(&cop.Spec, spec, strategyFunc); err != nil {
		_ = fmt.Errorf("ERROR: %s", err)
	}
	if reflect.DeepEqual(cop, node) {
		fmt.Printf("Node unchanged")
	}
	r.Client.Update(context.Background(), cop, &client.UpdateOptions{})
}

func (r *NodeLabelerReconciler) ManageNodes(nodes *corev1.NodeList, nodeLabelerSpec v1alpha1.NodeLabelerSpec) {
	for _, node := range nodes.Items {
		r.AssignAttributesToNodes(&node, nodeLabelerSpec.Merge.ObjectMeta, nodeLabelerSpec.Merge.NodeSpec, mergo.WithAppendSlice)
		// r.AssignAttributesToNodes(&node, nodeLabelerSpec.Overwrite.ObjectMeta, nodeLabelerSpec.Overwrite.NodeSpec, mergo.WithOverride)
	}

}

func (r *NodeLabelerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	nodeLabeler := &kubebuilderv1alpha1.NodeLabeler{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, nodeLabeler)
	allNodes := r.getAllNodes(ctx)
	filteredNodes := corev1.NodeList{}
	expressions := nodeLabeler.Spec.NodeSelectorTerms
	for _, node := range allNodes.Items {
		if match := helpers.NodeMatchesNodeSelectorTerms(&node, expressions); match {
			filteredNodes.Items = append(filteredNodes.Items, node)
		}
	}
	r.ManageNodes(&filteredNodes, nodeLabeler.Spec)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeLabelerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubebuilderv1alpha1.NodeLabeler{}).
		Complete(r)
}
