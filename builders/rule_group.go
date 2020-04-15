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
	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
)

// RuleGroupBuilder provides an interface to build RuleGroup
type RuleGroupBuilder struct {
	obj   *monitoringv1.RuleGroup
	rules []*RuleBuilder
}

// RuleGroup returns a new rulegroup builder
func RuleGroup() *RuleGroupBuilder {
	RuleGroup := &monitoringv1.RuleGroup{}
	return &RuleGroupBuilder{
		obj: RuleGroup,
	}
}

func (r *RuleGroupBuilder) Name(Name string) *RuleGroupBuilder {
	r.obj.Name = Name
	return r
}

func (r *RuleGroupBuilder) Interval(Interval string) *RuleGroupBuilder {
	r.obj.Interval = Interval
	return r
}

func (r *RuleGroupBuilder) Rules(Rules ...*RuleBuilder) *RuleGroupBuilder {
	r.rules = Rules
	return r
}

func (r *RuleGroupBuilder) PartialResponseStrategy(prs string) *RuleGroupBuilder {
	r.obj.PartialResponseStrategy = prs
	return r
}

// Build returns the object after making certain assertions
func (r *RuleGroupBuilder) Build() (monitoringv1.RuleGroup, error) {

	r.obj.Rules = []monitoringv1.Rule{}
	for _, rBuilder := range r.rules {
		rule, err := rBuilder.Build()
		if err != nil {
			return monitoringv1.RuleGroup{}, err
		}

		r.obj.Rules = append(r.obj.Rules, rule)
	}

	return *r.obj, nil
}
