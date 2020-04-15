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
	"errors"

	"github.com/alecthomas/units"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ContainerBuilder provides an interface to build containers
type ContainerBuilder struct {
	obj             *corev1.Container
	securityContext *SecurityContextBuilder
}

// Container returns a new container builder
func Container(name string, image string) *ContainerBuilder {
	container := &corev1.Container{
		Name:                     name,
		Image:                    image,
		ImagePullPolicy:          corev1.PullAlways,
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
	}

	return &ContainerBuilder{
		obj: container,
	}
}

// Args sets the arguments for that container
func (c *ContainerBuilder) Args(args ...string) *ContainerBuilder {
	c.obj.Args = args
	return c
}

// SecurityContext sets the SecurityContext for that container
func (c *ContainerBuilder) SecurityContext(SecurityContext *SecurityContextBuilder) *ContainerBuilder {
	c.securityContext = SecurityContext
	return c
}

// Port appends a port to the container
func (c *ContainerBuilder) Port(name string, port int32) *ContainerBuilder {
	c.obj.Ports = append(c.obj.Ports, v1.ContainerPort{
		Name:          name,
		ContainerPort: port,
		Protocol:      corev1.ProtocolTCP,
	})
	return c
}

// Volume appends a volume to the container
func (c *ContainerBuilder) Volume(name string, path string) *ContainerBuilder {
	c.obj.VolumeMounts = append(c.obj.VolumeMounts, v1.VolumeMount{
		Name:      name,
		MountPath: path,
	})
	return c
}

// Resources defines the resource configuration for the container
func (c *ContainerBuilder) Resources(cpu int64, memory int64, storage int64, factor float64) *ContainerBuilder {
	memory = memory * int64(units.Mebibyte)
	storage = storage * int64(units.Megabyte)

	cpuLimit := int64(float64(cpu) * factor)
	memoryLimit := int64(float64(memory) * factor)
	storageLimit := int64(float64(storage) * factor)

	c.obj.Resources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:              *resource.NewMilliQuantity(cpuLimit, resource.DecimalSI),
			v1.ResourceMemory:           *resource.NewQuantity(memoryLimit, resource.BinarySI),
			v1.ResourceEphemeralStorage: *resource.NewQuantity(storageLimit, resource.DecimalSI),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:              *resource.NewMilliQuantity(cpu, resource.DecimalSI),
			v1.ResourceMemory:           *resource.NewQuantity(memory, resource.BinarySI),
			v1.ResourceEphemeralStorage: *resource.NewQuantity(storage, resource.DecimalSI),
		},
	}

	return c
}

// HTTPProbe creates both a readiness and liveness probe with provided intervals
func (c *ContainerBuilder) HTTPProbe(port string, path string, readyInterval int32, liveInterval int32) *ContainerBuilder {
	handler := v1.Handler{
		HTTPGet: &v1.HTTPGetAction{
			Path:   path,
			Port:   intstr.FromString(port),
			Scheme: v1.URISchemeHTTP,
		},
	}

	return c.Probe(handler, readyInterval, liveInterval)
}

// PortProbe creates both a readiness and liveness probe with provided intervals
func (c *ContainerBuilder) PortProbe(port string, readyInterval int32, liveInterval int32) *ContainerBuilder {
	handler := v1.Handler{
		TCPSocket: &v1.TCPSocketAction{
			Port: intstr.FromString(port),
		},
	}

	return c.Probe(handler, readyInterval, liveInterval)
}

// Probe creates both a readiness and liveness probe based on a handler provided
func (c *ContainerBuilder) Probe(handler v1.Handler, readyInterval int32, liveInterval int32) *ContainerBuilder {
	c.obj.ReadinessProbe = &v1.Probe{
		Handler:             handler,
		InitialDelaySeconds: 0,
		PeriodSeconds:       readyInterval,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	c.obj.LivenessProbe = &v1.Probe{
		Handler:             handler,
		InitialDelaySeconds: 0,
		PeriodSeconds:       liveInterval,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	return c
}

// EnvVarFromString register one environment variable set from the string pair.
func (c *ContainerBuilder) EnvVarFromString(name string, value string) *ContainerBuilder {
	c.obj.Env = append(c.obj.Env, corev1.EnvVar{
		Name:  name,
		Value: value,
	})
	return c
}

// EnvVarFromConfigMap register one environment variable set from the configMap.
func (c *ContainerBuilder) EnvVarFromConfigMap(name string, cfmName string, cfmKey string) *ContainerBuilder {
	c.obj.Env = append(c.obj.Env, corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cfmName,
				},
				Key: cfmKey,
			},
		},
	})
	return c
}

// EnvVarFromSecret register one environment variable set from the secret.
func (c *ContainerBuilder) EnvVarFromSecret(name string, scName string, scKey string) *ContainerBuilder {
	c.obj.Env = append(c.obj.Env, corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: scName,
				},
				Key: scKey,
			},
		},
	})
	return c
}

// Build returns the object after making certain assertions
func (c *ContainerBuilder) Build() (corev1.Container, error) {
	if c.securityContext == nil {
		return corev1.Container{}, errors.New("missing security context")
	}
	securityContext, err := c.securityContext.Build()

	if err != nil {
		return corev1.Container{}, err
	}
	c.obj.SecurityContext = &securityContext
	return *c.obj, nil
}
