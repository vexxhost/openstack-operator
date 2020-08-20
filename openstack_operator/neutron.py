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
Neutron service

This code takes care of doing the operations of the OpenStack Neutron API
service.
"""

from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(spec, **_):
    """Create and re-sync a Neutron instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """
    utils.ensure_mysql_cluster("neutron", spec["mysql"])
    utils.deploy_rabbitmq("neutron")
    utils.create_or_update('neutron/daemonset-server.yml.j2', spec=spec)
    utils.create_or_update('neutron/daemonset-openvswitch-agent.yml.j2',
                           spec=spec)
    utils.create_or_update('neutron/daemonset-l3-agent.yml.j2', spec=spec)
    utils.create_or_update('neutron/daemonset-dhcp-agent.yml.j2', spec=spec)
    utils.create_or_update('neutron/daemonset-metadata-agent.yml.j2',
                           spec=spec)
    utils.create_or_update('neutron/service.yml.j2')

    identity.ensure_application_credential(name="neutron")

    url = None
    if "ingress" in spec:
        utils.create_or_update('neutron/ingress.yml.j2', spec=spec)
        url = spec["ingress"]["host"]

    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="neutron", service_type="network",
                                url=url, desc="Neutron Service")
