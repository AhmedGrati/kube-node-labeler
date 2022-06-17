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
	kubebuilderv1alpha1 "kube-node-labeler/api/v1alpha1"
	"kube-node-labeler/helpers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
	for key, node := range nodes.Items {
		fmt.Printf("Node Number %v: %s\n", key, node.Name)
	}
	return *nodes
}

func (r *NodeLabelerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	nodeLabeler := &kubebuilderv1alpha1.NodeLabeler{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, nodeLabeler)
	_ = r.getAllNodes(ctx)
	expressions := nodeLabeler.Spec.NodeSelectorTerms
	for _, expression := range expressions {
		matchExpressions := expression.MatchExpressions
		matchFields := expression.MatchFields

		if len(matchExpressions) > 0 {
			_ = helpers.HandleMatchExpressions(matchExpressions)
			_ = helpers.HandleMatchFields(matchFields)
		}

	}
	l.Info("The Node Labeler is ready to work", "spec", nodeLabeler.Spec)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeLabelerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubebuilderv1alpha1.NodeLabeler{}).
		Complete(r)
}
