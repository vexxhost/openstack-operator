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

"""cinder Operator

This module maintains the operator for Cinder.
"""

from openstack_operator import database
from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync a cinder instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    # deploy mysql for cinder
    database.ensure_mysql_cluster("cinder", spec=spec["mysql"])

    # deploy rabbitmq
    utils.deploy_rabbitmq("cinder")

    # deploy cinder
    config_hash = utils.generate_hash(spec)

    for component in ("api", "scheduler", "volume"):
        utils.create_or_update('cinder/daemonset.yml.j2',
                               name=name, spec=spec,
                               component=component,
                               config_hash=config_hash)

    utils.create_or_update('cinder/service.yml.j2', name=name)

    url = None
    if "ingress" in spec:
        utils.create_or_update('cinder/ingress.yml.j2',
                               name=name, spec=spec)

        url = spec["ingress"]["host"]

    # Create application credential
    identity.ensure_application_credential(name="cinder")

    identity.ensure_service(name="cinder", service_type="block-storage",
                            url=url, desc="Cinder Volume Service",
                            path="/v3/$(project_id)s")
    identity.ensure_service(name="cinderv2", service_type="volumev2",
                            url=url, desc="Cinder Volume Service V2",
                            path="/v2/$(project_id)s", internal="cinder")
    identity.ensure_service(name="cinderv3", service_type="volumev3",
                            url=url, desc="Cinder Volume Service V3",
                            path="/v3/$(project_id)s", internal="cinder")


def update(name, spec, **_):
    """Update a cinder

    This function updates the deployment for cinder if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('cinder/ingress.yml.j2',
                               name=name, spec=spec)
