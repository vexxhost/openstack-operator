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


class MagnumAPIDeploymentTestCase(base.DeploymentTestCase):
    """Basic tests for the API Deployment."""

    RELEASE_TYPE = 'magnum'
    TEMPLATE_FILE = 'magnum/deployment.yml.j2'
    TEMPLATE_PARAMS = {'component': 'api'}


class MagnumConductorDeploymentTestCase(base.DeploymentTestCase):
    """Basic tests for the Conductor Deployment."""

    RELEASE_TYPE = 'magnum'
    TEMPLATE_FILE = 'magnum/deployment.yml.j2'
    TEMPLATE_PARAMS = {'component': 'conductor'}
    PORT_EXPOSED = False
