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

from openstack_operator import database
from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync any Magnum instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    if "mysql" not in spec:
        spec["mysql"] = {}
    database.ensure_mysql_cluster("magnum", spec["mysql"])

    # deploy rabbitmq
    if not utils.ensure_secret("openstack", "magnum-rabbitmq"):
        utils.create_or_update('magnum/secret-rabbitmq.yml.j2',
                               password=utils.generate_password())
    utils.create_or_update('magnum/rabbitmq.yml.j2', spec=spec)

    # deploy magnum
    config_hash = utils.generate_hash(spec)
    for component in ("api", "conductor"):
        utils.create_or_update('magnum/daemonset.yml.j2',
                               name=name, spec=spec,
                               component=component,
                               config_hash=config_hash)

    utils.create_or_update('magnum/service.yml.j2',
                           name=name)

    url = None
    if "ingress" in spec:
        utils.create_or_update('magnum/ingress.yml.j2',
                               name=name, spec=spec)
        url = spec["ingress"]["host"]

    # Create application credential
    identity.ensure_application_credential(name="magnum")

    # Create service and endpoints
    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="magnum-api", path="/v1",
                                service_type="container-infra", url=url,
                                desc="Container Infrastructure \
                                    Management Service")


def update(name, spec, **_):
    """Update a Magnum

    This function updates the deployment for Magnum if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('magnum/ingress.yml.j2',
                               name=name, spec=spec)
