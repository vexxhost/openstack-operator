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

"""horizon Operator

This module maintains the operator for Mcrouter, it takes care of creating
the appropriate deployments, Mcrouter, pod monitors and Prometheus rules.
"""

import kopf

from openstack_operator import utils


@kopf.on.create('dashboard.openstack.org', 'v1alpha1', 'horizons')
def create_secret(name, **_):
    """Create a new horizon secret"""

    res = utils.get_secret("openstack", name)
    if res is None:
        utils.create_or_update('horizon/secret-secretkey.yml.j2',
                               name=name,
                               secret=utils.generate_password())


@kopf.on.resume('dashboard.openstack.org', 'v1alpha1', 'horizons')
@kopf.on.create('dashboard.openstack.org', 'v1alpha1', 'horizons')
def create_or_resume(name, spec, **_):
    """Create and re-sync a horizon instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    # Grab the secretkey secret
    conn = utils.get_openstack_connection()
    auth_url = conn.config.auth["auth_url"]
    config = utils.create_or_update('horizon/configmap.yml.j2',
                                    name=name, spec=spec, auth_url=auth_url)
    config_hash = utils.generate_hash(config.obj['data'])
    env = utils.get_uwsgi_env()
    utils.create_or_update('horizon/deployment.yml.j2',
                           config_hash=config_hash, name=name,
                           spec=spec, env=env)
    utils.create_or_update('horizon/service.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('horizon/memcached.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('horizon/horizontalpodautoscaler.yml.j2',
                           name=name)
    if "ingress" in spec:
        utils.create_or_update('horizon/ingress.yml.j2',
                               name=name, spec=spec)


@kopf.on.update('dashboard.openstack.org', 'v1alpha1', 'horizons')
def update(name, spec, **_):
    """Update a horizon

    This function updates the deployment for horizon if there are any
    changes that happen within it.
    """
    conn = utils.get_openstack_connection()
    auth_url = conn.config.auth["auth_url"]
    config = utils.create_or_update('horizon/configmap.yml.j2',
                                    name=name, spec=spec, auth_url=auth_url)
    config_hash = utils.generate_hash(config.obj['data'])
    env = utils.get_uwsgi_env()
    utils.create_or_update('horizon/deployment.yml.j2',
                           config_hash=config_hash, name=name,
                           spec=spec, env=env)
    if "ingress" in spec:
        utils.create_or_update('horizon/ingress.yml.j2',
                               name=name, spec=spec)
