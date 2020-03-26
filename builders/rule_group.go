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
