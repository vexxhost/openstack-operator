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

"""Tests for Keystone Operator

This module contains all the tests for the Keystone operator.
"""

from openstack_operator.tests.unit import base


class KeystoneDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the Daemonset."""

    RELEASE_TYPE = 'keystone'
    TEMPLATE_FILE = 'keystone/daemonset.yml.j2'


class KeystoneIngressTestCase(base.IngressTestCase):
    """Basic tests for the Ingress."""

    RELEASE_TYPE = 'keystone'
    TEMPLATE_FILE = 'keystone/ingress.yml.j2'


class KeystoneServiceTestCase(base.ServiceTestCase):
    """Basic tests for the Service."""

    RELEASE_TYPE = 'keystone'
    TEMPLATE_FILE = 'keystone/service.yml.j2'
