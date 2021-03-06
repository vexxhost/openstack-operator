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

"""Tests for Magnum Operator

This module contains all the tests for the Magnum operator.
"""

from openstack_operator.tests.unit import base


class MagnumAPIDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the API Daemonset."""

    RELEASE_TYPE = 'magnum'
    TEMPLATE_FILE = 'magnum/api/daemonset.yml.j2'


class MagnumConductorDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the Conductor Daemonset."""

    RELEASE_TYPE = 'magnum'
    TEMPLATE_FILE = 'magnum/conductor/daemonset.yml.j2'
    PORT_EXPOSED = False


class MagnumAPIServiceTestCase(base.ServiceTestCase):
    """Basic tests for the Service."""

    RELEASE_TYPE = 'magnum'
    TEMPLATE_FILE = 'magnum/api/service.yml.j2'
