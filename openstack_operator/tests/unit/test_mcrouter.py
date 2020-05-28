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

"""Tests for Mcrouter Operator

This module contains all the tests for the Mcrouter operator.
"""

from openstack_operator.tests.unit import base


class McrouterDeploymentTestCase(base.StatefulSetTestCase):
    """Basic tests for the Deployment."""

    SAMPLE_FILE = 'infrastructure_v1alpha1_mcrouter.yaml'
    TEMPLATE_FILE = 'mcrouter/deployment.yml.j2'
    AUTO_GENERATED = False
