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

"""Tests for Heat Operator

This module contains all the tests for the Heat operator.
"""

from openstack_operator.tests.unit import base


class HeatAPIDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the api Daemonset."""

    RELEASE_TYPE = 'heat'
    TEMPLATE_FILE = 'heat/api/daemonset.yml.j2'
    PORT_EXPOSED = False

    def test_envvar_default_host_exists(self):
        """Ensure that heat api daemonset has OS_DEFAULT__HOST env var
        to set the engine host"""
        envvar_name_list = []
        envvar_list = \
            self.object['spec']['template']['spec']['containers'][0]["env"]
        for envvar in envvar_list:
            envvar_name_list.append(envvar["name"])
        self.assertIn('OS_DEFAULT__HOST', envvar_name_list)


class HeatAPIServiceTestCase(base.ServiceTestCase):
    """Basic tests for the api Service."""

    RELEASE_TYPE = 'heat'
    TEMPLATE_FILE = 'heat/api/service.yml.j2'


class HeatAPIcfnDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the api cfn Daemonset."""

    RELEASE_TYPE = 'heat'
    TEMPLATE_FILE = 'heat/api-cfn/daemonset.yml.j2'
    PORT_EXPOSED = False

    def test_envvar_default_host_exists(self):
        """Ensure that heat api cfn daemonset has OS_DEFAULT__HOST env var
        to set the engine host"""
        envvar_name_list = []
        envvar_list = \
            self.object['spec']['template']['spec']['containers'][0]["env"]
        for envvar in envvar_list:
            envvar_name_list.append(envvar["name"])
        self.assertIn('OS_DEFAULT__HOST', envvar_name_list)


class HeatAPIcfnServiceTestCase(base.ServiceTestCase):
    """Basic tests for the api cfn Service."""

    RELEASE_TYPE = 'heat'
    TEMPLATE_FILE = 'heat/api-cfn/service.yml.j2'


class HeatEngineDaemonsetTestCase(base.DaemonSetTestCase):
    """Basic tests for the engine Daemonset."""

    RELEASE_TYPE = 'heat'
    TEMPLATE_FILE = 'heat/engine/daemonset.yml.j2'
    PORT_EXPOSED = False

    def test_envvar_default_host_exists(self):
        """Ensure that heat daemonset has OS_DEFAULT__HOST env var
        to set the engine host"""
        envvar_name_list = []
        envvar_list = \
            self.object['spec']['template']['spec']['containers'][0]["env"]
        for envvar in envvar_list:
            envvar_name_list.append(envvar["name"])
        self.assertIn('OS_DEFAULT__HOST', envvar_name_list)
