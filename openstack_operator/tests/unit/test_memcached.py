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

from unittest import mock

from openstack_operator import memcached
from openstack_operator.tests.unit import base


class MemcachedOperatorTestCase(base.BaseTestCase):
    """Basic tests for the operator."""

    @mock.patch.object(memcached.utils, 'create_or_update')
    @mock.patch.object(memcached.utils, 'ensure_absent')
    def test_ensure_deployment_removal(self, mock_ensure_absent, _):
        """Test that we remove the old deployment"""
        memcached.create_or_resume("foo", {})
        mock_ensure_absent.assert_called_once_with(
            'memcached/deployment.yml.j2', name="foo", spec={})


class MemcachedStatefulSetTestCase(base.StatefulSetTestCase):
    """Basic tests for the StatefulSet."""

    SAMPLE_FILE = 'infrastructure_v1alpha1_memcached.yaml'
    TEMPLATE_FILE = 'memcached/statefulset.yml.j2'
    AUTO_GENERATED = False
