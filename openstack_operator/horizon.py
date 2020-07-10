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


from openstack_operator import utils


def create_secret(name, **_):
    """Create a new horizon secret"""

    res = utils.get_secret("openstack", name)
    if res is None:
        utils.create_or_update('horizon/secret-secretkey.yml.j2',
                               name=name,
                               secret=utils.generate_password())


def create_or_resume(name, spec, **_):
    """Create and re-sync a horizon instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    # Grab the secretkey secret
    config = utils.create_or_update('horizon/configmap.yml.j2',
                                    name=name, spec=spec)
    config_hash = utils.generate_hash(config.obj['data'])

    utils.create_or_update('horizon/daemonset.yml.j2',
                           config_hash=config_hash, name=name,
                           spec=spec)
    utils.create_or_update('horizon/service.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('horizon/memcached.yml.j2',
                           name=name, spec=spec)
    if "ingress" in spec:
        utils.create_or_update('horizon/ingress.yml.j2',
                               name=name, spec=spec)

    # NOTE(Alex): We should remove this once all deployments are no longer
    #               using Deployment.
    utils.ensure_absent('horizon/deployment.yml.j2',
                        config_hash=config_hash, name=name,
                        spec=spec)

    # NOTE(Alex): We should remove this once all deployments are no longer
    #               using HPA.
    utils.ensure_absent('horizon/horizontalpodautoscaler.yml.j2',
                        name=name)


def update(name, spec, **_):
    """Update a horizon

    This function updates the deployment for horizon if there are any
    changes that happen within it.
    """
    config = utils.create_or_update('horizon/configmap.yml.j2',
                                    name=name, spec=spec)
    config_hash = utils.generate_hash(config.obj['data'])

    utils.create_or_update('horizon/daemonset.yml.j2',
                           config_hash=config_hash, name=name,
                           spec=spec)
    if "ingress" in spec:
        utils.create_or_update('horizon/ingress.yml.j2',
                               name=name, spec=spec)
