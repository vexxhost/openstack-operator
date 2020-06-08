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

"""Magnum Operator

This module maintains the operator for Magnum, it takes care of creating
the appropriate deployments, an instance of Memcache, RabbitMQ and a database
server for the installation.
"""

from openstack_operator import utils


def create_or_resume(name, spec, **_):
    """Create and re-sync any Magnum instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    env = utils.get_uwsgi_env()
    config_hash = utils.generate_hash(spec)
    for component in ("api", "conductor"):
        utils.create_or_update('magnum/deployment.yml.j2',
                               name=name, spec=spec,
                               component=component, env=env,
                               config_hash=config_hash)
        utils.create_or_update('magnum/horizontalpodautoscaler.yml.j2',
                               name=name, component=component)

    utils.create_or_update('magnum/service.yml.j2',
                           name=name)

    if "ingress" in spec:
        utils.create_or_update('magnum/ingress.yml.j2',
                               name=name, spec=spec)


def update(name, spec, **_):
    """Update a Magnum

    This function updates the deployment for Magnum if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('magnum/ingress.yml.j2',
                               name=name, spec=spec)
