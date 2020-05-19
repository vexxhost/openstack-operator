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

"""Keystone Operator

This module maintains the operator for Keystone which does everything from
deployment to taking care of rotating fernet & credentials keys."""

import kopf

from openstack_operator import utils


@kopf.on.resume('identity.openstack.org', 'v1alpha1', 'keystones')
@kopf.on.create('identity.openstack.org', 'v1alpha1', 'keystones')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Keystone instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """
    env = utils.get_uwsgi_env()
    config_hash = utils.generate_hash(spec)
    utils.create_or_update('keystone/deployment.yml.j2',
                           name=name, spec=spec,
                           env=env, config_hash=config_hash)
    utils.create_or_update('keystone/service.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('keystone/horizontalpodautoscaler.yml.j2',
                           name=name)
    if "ingress" in spec:
        utils.create_or_update('keystone/ingress.yml.j2',
                               spec=spec)


@kopf.on.update('identity.openstack.org', 'v1alpha1', 'keystones')
def update(spec, **_):
    """Update a keystone

    This function updates the deployment for horizon if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('keystone/ingress.yml.j2',
                               spec=spec)
