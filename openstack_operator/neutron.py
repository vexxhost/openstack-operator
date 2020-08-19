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

from openstack_operator import database
from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(spec, **_):
    """Create and re-sync a Neutron instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    if not utils.ensure_secret("openstack", "neutron-rabbitmq"):
        utils.create_or_update('neutron/secret-rabbitmq.yml.j2',
                               password=utils.generate_password())

    database.ensure_mysql_cluster("neutron")

    utils.create_or_update('neutron/rabbitmq.yml.j2')
    utils.create_or_update('neutron/daemonset.yml.j2', spec=spec)
    utils.create_or_update('neutron/service.yml.j2')

    identity.ensure_application_credential(name="neutron")
