// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// PrometheusRuleBuilder provides an interface to build PrometheusRules
type PrometheusRuleBuilder struct {
	obj        *monitoringv1.PrometheusRule
	ruleGroups []*RuleGroupBuilder
	owner      metav1.Object
	scheme     *runtime.Scheme
}

// PrometheusRule returns a new PrometheusRule builder
func PrometheusRule(existing *monitoringv1.PrometheusRule, owner metav1.Object, scheme *runtime.Scheme) *PrometheusRuleBuilder {
	return &PrometheusRuleBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

func (pm *PrometheusRuleBuilder) Labels(labels map[string]string) *PrometheusRuleBuilder {
	pm.obj.Labels = labels
	return pm
}

// RuleGroups returns the ruleGroups
func (pm *PrometheusRuleBuilder) RuleGroups(ruleGroups ...*RuleGroupBuilder) *PrometheusRuleBuilder {
	pm.ruleGroups = ruleGroups
	return pm
}

// Build returns the object after making certain assertions
func (pm *PrometheusRuleBuilder) Build() error {
	pm.obj.Spec.Groups = []monitoringv1.RuleGroup{}
	for _, rgBuilder := range pm.ruleGroups {
		ruleGroup, err := rgBuilder.Build()
		if err != nil {
			return err
		}

		pm.obj.Spec.Groups = append(pm.obj.Spec.Groups, ruleGroup)
	}
	if !pm.isOwnedByOthers() {
		return controllerutil.SetControllerReference(pm.owner, pm.obj, pm.scheme)
	}
	return nil
}

// isOwnedByOthers checks if this podMonitor has been possessed by an another object already.
func (pm *PrometheusRuleBuilder) isOwnedByOthers() bool {
	ownerName := pm.owner.GetName()

	existingRefs := pm.obj.GetOwnerReferences()
	for _, r := range existingRefs {
		if r.Name == ownerName {
			return false
		} else if r.Controller != nil && *r.Controller {
			return true
		}
	}
	return false
}
