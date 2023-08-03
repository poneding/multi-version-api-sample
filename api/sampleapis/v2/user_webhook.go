/*
Copyright 2023.

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

package v2

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var userlog = logf.Log.WithName("user-resource")

func (r *User) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-sampleapis-poneding-com-v2-user,mutating=true,failurePolicy=fail,groups=sampleapis.poneding.com,resources=users,verbs=create;update,versions=v2,name=muser.dp.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Defaulter = &User{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *User) Default() {

}

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-sampleapis-poneding-com-v2-user,mutating=false,failurePolicy=fail,groups=sampleapis.poneding.com,resources=users,versions=v2,name=vuser.dp.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Validator = &User{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *User) ValidateCreate() (warnings admission.Warnings, err error) {
	userlog.Info("validate create", "name", r.Name)

	return make(admission.Warnings, 0), r.validateUser()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *User) ValidateUpdate(old runtime.Object) (warnings admission.Warnings, err error) {
	userlog.Info("validate update", "name", r.Name)

	return make(admission.Warnings, 0), r.validateUser()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *User) ValidateDelete() (warnings admission.Warnings, err error) {
	userlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return make(admission.Warnings, 0), nil
}

func (r *User) validateUser() error {
	var allErrs field.ErrorList
	if err := r.validateUserName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateUserSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "sampleapis.poneding.com", Kind: "User"},
		r.Name, allErrs)
}

func (r *User) validateUserSpec() *field.Error {
	return nil
}

func (r *User) validateUserName() *field.Error {
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
