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

"""Memcached Operator

This module maintains the operator for Memcached, it takes care of creating
the appropriate deployments, Mcrouter, pod monitors and Prometheus rules.
"""

import kopf

from openstack_operator import utils


@kopf.on.resume('infrastructure.vexxhost.cloud', 'v1alpha1', 'memcacheds')
@kopf.on.create('infrastructure.vexxhost.cloud', 'v1alpha1', 'memcacheds')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Memcached instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    utils.create_or_update('memcached/statefulset.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('memcached/service.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('memcached/mcrouter.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('memcached/podmonitor.yml.j2',
                           name=name, spec=spec, adopt=True)
    utils.create_or_update('memcached/prometheusrule.yml.j2',
                           name=name, spec=spec, adopt=True)

    # NOTE(mnaser): We should remove this once all deployments are no longer
    #               using Deployment for Memcached.
    utils.ensure_absent('memcached/deployment.yml.j2',
                        name=name, spec=spec)


@kopf.on.update('infrastructure.vexxhost.cloud', 'v1alpha1', 'memcacheds')
def update(name, spec, **_):
    """Update a Memcached

    This function updates the deployment for Memcached if there are any
    changes that happen within it.
    """

    utils.create_or_update('memcached/statefulset.yml.j2',
                           name=name, spec=spec, adopt=True)
