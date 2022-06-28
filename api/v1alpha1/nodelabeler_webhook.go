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

package v1alpha1

import (
	"regexp"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	// "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var nodelabelerlog = logf.Log.WithName("nodelabeler-resource")

func (r *NodeLabeler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-kubebuilder-kube-node-labeler-io-v1alpha1-nodelabeler,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubebuilder.kube.node.labeler.io,resources=nodelabelers,verbs=create;update,versions=v1alpha1,name=mnodelabeler.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &NodeLabeler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NodeLabeler) Default() {
	nodelabelerlog.Info("default", "name", r.Name)
	if r.Spec.Size == nil {
		r.Spec.Size = new(int)
		*r.Spec.Size = 0
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-kubebuilder-kube-node-labeler-io-v1alpha1-nodelabeler,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubebuilder.kube.node.labeler.io,resources=nodelabelers,verbs=create;update,versions=v1alpha1,name=vnodelabeler.kb.io,admissionReviewVersions=v1

// var _ webhook.Validator = &NodeLabeler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NodeLabeler) ValidateCreate() error {
	nodelabelerlog.Info("validate create", "name", r.Name)

	return r.validateNodeLabeler()
}

func (r *NodeLabeler) validateNodeLabeler() error {
	var allErrors field.ErrorList
	if err := r.validateNodeLabelerName(); err != nil {
		allErrors = append(allErrors, err)
	}
	if err := r.validateNodeLabelerSpec(); err != nil {
		allErrors = append(allErrors, err)
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "kubebuilder.kube.node.labeler.io", Kind: "NodeLabeler"},
		r.Name, allErrors,
	)
}

func (r *NodeLabeler) validateNodeLabelerSpec() *field.Error {
	if *r.Spec.Size < 0 {
		return field.Invalid(field.NewPath("spec").Child("size"), r.Spec.Size, "must be equal or greater than 0")
	}
	for _, regex := range r.Spec.NodeNamePatterns {

		if _, err := regexp.Compile(regex); err != nil {
			return field.Invalid(field.NewPath("spec").Child("nodeNamePatterns"), r.Spec.NodeNamePatterns, "Must have valid regex")
		}
	}
	return nil
}

func (r *NodeLabeler) validateNodeLabelerName() *field.Error {

	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The cronjob controller appends
		// a 11-character suffix to the cronjob (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore cronjob
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}

// // ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
// func (r *NodeLabeler) ValidateUpdate(old runtime.Object) error {
// 	nodelabelerlog.Info("validate update", "name", r.Name)

// 	// TODO(user): fill in your validation logic upon object update.
// 	return nil
// }

// // ValidateDelete implements webhook.Validator so a webhook will be registered for the type
// func (r *NodeLabeler) ValidateDelete() error {
// 	nodelabelerlog.Info("validate delete", "name", r.Name)

// 	// TODO(user): fill in your validation logic upon object deletion.
// 	return nil
// }
