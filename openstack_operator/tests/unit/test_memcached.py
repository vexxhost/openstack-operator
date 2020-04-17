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

"""Tests for Memcached Operator

This module contains all the tests for the Memcached operator.
"""

# Disable no-self-use
# pylint: disable=R0201

import mock

from oslotest import base

from openstack_operator import memcached


class MemcachedListTestCase(base.BaseTestCase):
    """Tests for determining server list."""

    @mock.patch.object(memcached.utils, 'get_ready_pod_ips')
    @mock.patch.object(memcached.utils, 'create_or_update')
    def test_with_no_ips(self, mock_create, mock_get_ready_pods):
        """Test a deployment with no ready pods."""

        mock_get_ready_pods.return_value = []
        memcached.deployment_event("default", {}, {})

        mock_create.assert_called_once_with('memcached/mcrouter.yml.j2',
                                            name=None, servers=[], spec={})

    @mock.patch.object(memcached.utils, 'get_ready_pod_ips')
    @mock.patch.object(memcached.utils, 'create_or_update')
    def test_with_single_ip(self, mock_create, mock_get_ready_pods):
        """Test a deployment with a single ready pod."""

        mock_get_ready_pods.return_value = ['1.1.1.1']
        memcached.deployment_event("default", {}, {})

        mock_create.assert_called_once_with(
            'memcached/mcrouter.yml.j2', name=None,
            servers=['1.1.1.1:11211'], spec={})

    @mock.patch.object(memcached.utils, 'get_ready_pod_ips')
    @mock.patch.object(memcached.utils, 'create_or_update')
    def test_multiple_ips(self, mock_create, mock_get_ready_pods):
        """Test a deployment with a multiple ready pods."""

        mock_get_ready_pods.return_value = ['1.1.1.1', '2.2.2.2']
        memcached.deployment_event("default", {}, {})

        mock_create.assert_called_once_with(
            'memcached/mcrouter.yml.j2', name=None,
            servers=['1.1.1.1:11211', '2.2.2.2:11211'], spec={})
