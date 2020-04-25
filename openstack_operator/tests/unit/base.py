# Copyright 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Base test classes

This module contains the base test classes.
"""

import testtools

import yaml

from openstack_operator import utils


class BaseTestCase(testtools.TestCase):
    """Base test class for the OpenStack operator."""


class KubernetesObjectTestCase(testtools.TestCase):
    """Base class for Kubernetes object tests."""

    SAMPLES_PATH = 'config/samples'

    @classmethod
    def setUpClass(cls):
        sample_path = "%s/%s" % (cls.SAMPLES_PATH, cls.SAMPLE_FILE)
        with open(sample_path) as sample_fd:
            sample = yaml.load(sample_fd, Loader=yaml.FullLoader)
        name = sample['metadata']['name']
        spec = sample['spec']

        cls.object = utils.render_template(cls.TEMPLATE_FILE,
                                           name=name, spec=spec)


class KubernetesAppTestCaseMixin:
    """Mix-in to be used for tests that involve apps and containers."""

    def test_containers_use_always_image_pull_policy(self):
        """Ensure that all containers use 'Always' as imagePullPolicy."""
        for container in self.object['spec']['template']['spec']['containers']:
            self.assertEqual("Always", container.get('imagePullPolicy'))

    def test_containers_have_liveness_probe(self):
        """Ensure that all containers have liveness probes."""
        for container in self.object['spec']['template']['spec']['containers']:
            self.assertIn('livenessProbe', container)

    def test_containers_have_readiness_probe(self):
        """Ensure that all containers have readiness probes."""
        for container in self.object['spec']['template']['spec']['containers']:
            self.assertIn('readinessProbe', container)

    def test_containers_have_resource_limits(self):
        """Ensure that all containers have resource limits."""
        for container in self.object['spec']['template']['spec']['containers']:
            self.assertIn('resources', container)

    def test_container_http_probes_have_no_metrics_path(self):
        """Ensure that http probes (liveness/rediness) of all containers
         don't have metrics path"""
        for container in self.object['spec']['template']['spec']['containers']:
            if 'httpGet' in container['readinessProbe']:
                self.assertNotEqual(
                    container['readinessProbe']['httpGet']['path'],
                    '/metrics'
                )
            if 'httpGet' in container['livenessProbe']:
                self.assertNotEqual(
                    container['livenessProbe']['httpGet']['path'],
                    '/metrics'
                )


class DeploymentTestCase(KubernetesObjectTestCase,
                         KubernetesAppTestCaseMixin):
    """Basic tests for Kubernetes Deployments."""


class StatefulSetTestCase(KubernetesObjectTestCase,
                          KubernetesAppTestCaseMixin):
    """Basic tests for Kubernetes StatefulSets."""
