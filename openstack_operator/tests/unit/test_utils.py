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

"""Tests for utilities

This module contains all the tests for the utilities
"""

import testtools

from openstack_operator import utils


class LabelsTestCase(testtools.TestCase):
    """Base test class for the OpenStack operator."""

    def test_labels(self):
        """Test basic label generation."""
        labels = utils.labels("foo", "bar")

        self.assertEqual(labels, """app.kubernetes.io/instance: bar
app.kubernetes.io/managed-by: openstack-operator
app.kubernetes.io/name: foo""")

    def test_labels_with_component(self):
        """Test label generation with components."""
        labels = utils.labels("foo", "bar", "baz")

        self.assertEqual(labels, """app.kubernetes.io/component: baz
app.kubernetes.io/instance: bar
app.kubernetes.io/managed-by: openstack-operator
app.kubernetes.io/name: foo""")
