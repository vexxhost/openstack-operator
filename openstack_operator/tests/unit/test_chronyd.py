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

"""Tests for Chronyd Operator

This module contains all the tests for the Chronyd operator.
"""

from openstack_operator.tests.unit import base


class ChronydAPIDeploymentTestCase(base.DaemonSetTestCase):
    """Basic tests for the DaemonSet."""

    RELEASE_TYPE = 'chronyd'
    TEMPLATE_FILE = 'chronyd/daemonset.yml.j2'
