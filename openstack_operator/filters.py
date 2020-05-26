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

"""Kopf filters

This module contains a few common filters to be used throughout the operator
in order to reduce strain on the API server.
"""


def managed(namespace, labels, **_):
    """Check if a resource is managed by the operator."""
    return namespace == 'openstack' and \
        labels.get('app.kubernetes.io/managed-by') == 'openstack-operator'
