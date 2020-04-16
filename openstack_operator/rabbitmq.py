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

"""Rabbitmq Operator

This module maintains the operator for Rabbitmq, it takes care of creating
the appropriate deployments, Rabbitmq, pod monitors and Prometheus rules.
"""

import kopf

from openstack_operator import utils


@kopf.on.resume('infrastructure.vexxhost.cloud', 'v1alpha1', 'rabbitmqs')
@kopf.on.create('infrastructure.vexxhost.cloud', 'v1alpha1', 'rabbitmqs')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Rabbitmq instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    utils.create_or_update('rabbitmq/deployment.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('rabbitmq/service.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('rabbitmq/podmonitor.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('rabbitmq/prometheusrule.yml.j2',
                           name=name, spec=spec)


@kopf.on.update('infrastructure.vexxhost.cloud', 'v1alpha1', 'rabbitmqs')
def update(name, spec, **_):
    """Update a Rabbitmq

    This function updates the deployment for Rabbitmq if there are any
    changes that happen within it.
    """

    utils.create_or_update('rabbitmq/deployment.yml.j2',
                           name=name, spec=spec)
