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

"""
Nova service

This code takes care of doing the operations of the OpenStack Nova API
service.
"""

import kopf

from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True

# NOTE(mnaser): Implement dynamic cells
CELLS = [
    'cell0',
    'cell1'
]


def create_or_resume(spec, **_):
    """Create and re-sync a Nova instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    databases = {}

    identity.ensure_application_credential(name="nova")

    databases['api'] = utils.ensure_mysql_cluster(
        "nova-api", database="nova_api"
    )

    for cell in CELLS:
        databases[cell] = utils.ensure_mysql_cluster(
            "nova-%s" % cell, database="nova_%s" % cell)

        # NOTE(mnaser): cell0 does not need a message queue
        if cell != 'cell0':
            utils.deploy_rabbitmq("nova-%s" % cell)

    utils.create_or_update('nova/conductor/daemonset.yml.j2', spec=spec)
    utils.create_or_update('nova/scheduler/daemonset.yml.j2', spec=spec)

    utils.create_or_update('nova/metadata-api/daemonset.yml.j2', spec=spec)
    utils.create_or_update('nova/metadata-api/service.yml.j2')

    utils.create_or_update('nova/novncproxy/daemonset.yml.j2', spec=spec)
    utils.create_or_update('nova/novncproxy/service.yml.j2')

    utils.create_or_update('nova/compute-api/daemonset.yml.j2', spec=spec)
    utils.create_or_update('nova/compute-api/service.yml.j2')

    utils.create_or_update('nova/compute/daemonset.yml.j2', spec=spec)

    api_url = None
    if "ingress" in spec:
        utils.create_or_update('nova/ingress.yml.j2', spec=spec)
        api_url = spec["ingress"]["host"]["api"]

    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="nova",
                                service_type="compute",
                                url=api_url, path="/v2.1",
                                desc="OpenStack Compute")


@kopf.on.create('apps', 'v1', 'daemonsets', labels={
    'app.kubernetes.io/managed-by': 'openstack-operator',
    'app.kubernetes.io/name': 'nova',
    'app.kubernetes.io/component': 'conductor',
})
def run_database_migrations(**_):
    """Run database migrations

    This watches for any changes to the image ID for the Nova conductor
    deployment and triggers a database migrations
    """

    cell0 = utils.ensure_mysql_cluster("nova-cell0")
    utils.create_or_update('nova/conductor/job.yml.j2', adopt=True,
                           cell0_db=cell0['connection'])
