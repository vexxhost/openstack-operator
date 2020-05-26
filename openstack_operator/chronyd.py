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

"""chronyd Operator

This module maintains the operator for Chronyd, it takes care of creating
the appropriate daemonset, and configMap.
"""


from openstack_operator import utils


def create_or_resume(spec, **_):
    """Create and re-sync a chronyd instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    utils.create_or_update('chronyd/daemonset.yml.j2',
                           spec=spec)
    utils.create_or_update('chronyd/configmap.yml.j2',
                           spec=spec)
