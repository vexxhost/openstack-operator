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

"""Mcrouter Operator

This module maintains the operator for Mcrouter, it takes care of creating
the appropriate deployments, Mcrouter, pod monitors and Prometheus rules.
"""

import json
import kopf

from openstack_operator import utils


class McrouterSpecEncoder(json.JSONEncoder):
    """McrouterSpecEncoder makes kopf dictview class JSON serializable"""

    def default(self, o):                   # pylint: disable=E0202
        return o.__dict__["_src"]["spec"]


@kopf.on.resume('infrastructure.vexxhost.cloud', 'v1alpha1', 'mcrouters')
@kopf.on.create('infrastructure.vexxhost.cloud', 'v1alpha1', 'mcrouters')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Mcrouter instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """
    data = McrouterSpecEncoder().encode(spec)
    utils.create_or_update('mcrouter/configmap.yml.j2',
                           name=name, data=data, adopt=True)
    utils.create_or_update('mcrouter/deployment.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('mcrouter/service.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('mcrouter/podmonitor.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('mcrouter/prometheusrule.yml.j2',
                           name=name, spec=spec, adopt=True)


@kopf.on.update('infrastructure.vexxhost.cloud', 'v1alpha1', 'mcrouters')
def update(name, spec, **_):
    """Update a Mcrouter

    This function updates the deployment for Mcrouter if there are any
    changes that happen within it.
    """
    data = McrouterSpecEncoder().encode(spec)
    utils.create_or_update('mcrouter/configmap.yml.j2',
                           name=name, data=data, adopt=True)
    utils.create_or_update('mcrouter/deployment.yml.j2',
                           name=name, spec=spec, adopt=True)
