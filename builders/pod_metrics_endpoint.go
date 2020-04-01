package builders

import (
	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
)

// PodMetricsEndpointBuilder provides an interface to build podmonitors
type PodMetricsEndpointBuilder struct {
	obj *monitoringv1.PodMetricsEndpoint
}

// PodMonitor returns a new podmonitor builder
func PodMetricsEndpoint() *PodMetricsEndpointBuilder {
	podMetricsEndpoint := &monitoringv1.PodMetricsEndpoint{}
	return &PodMetricsEndpointBuilder{
		obj: podMetricsEndpoint,
	}
}

func (pme *PodMetricsEndpointBuilder) Port(port string) *PodMetricsEndpointBuilder {
	pme.obj.Port = port
	return pme
}

func (pme *PodMetricsEndpointBuilder) Path(path string) *PodMetricsEndpointBuilder {
	pme.obj.Path = path
	return pme
}

func (pme *PodMetricsEndpointBuilder) Scheme(scheme string) *PodMetricsEndpointBuilder {
	pme.obj.Scheme = scheme
	return pme
}

func (pme *PodMetricsEndpointBuilder) Params(params map[string][]string) *PodMetricsEndpointBuilder {
	pme.obj.Params = params
	return pme
}

func (pme *PodMetricsEndpointBuilder) Interval(interval string) *PodMetricsEndpointBuilder {
	pme.obj.Interval = interval
	return pme
}

func (pme *PodMetricsEndpointBuilder) ScrapeTimeout(scrapeTimeout string) *PodMetricsEndpointBuilder {
	pme.obj.ScrapeTimeout = scrapeTimeout
	return pme
}

func (pme *PodMetricsEndpointBuilder) HonorLabels(honorLabels bool) *PodMetricsEndpointBuilder {
	pme.obj.HonorLabels = honorLabels
	return pme
}

func (pme *PodMetricsEndpointBuilder) HonorTimestamps(honorTimestamps bool) *PodMetricsEndpointBuilder {
	pme.obj.HonorTimestamps = &honorTimestamps
	return pme
}

func (pme *PodMetricsEndpointBuilder) ProxyURL(proxyURL string) *PodMetricsEndpointBuilder {
	pme.obj.ProxyURL = &proxyURL
	return pme
}

// Build returns the object after making certain assertions
func (pme *PodMetricsEndpointBuilder) Build() (monitoringv1.PodMetricsEndpoint, error) {
	return *pme.obj, nil
}
