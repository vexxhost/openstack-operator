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

from openstack_operator import utils

MEMCACHED = True

# NOTE(mnaser): Implement dynamic cells
CELLS = [
    'cell0',
    'cell1'
]


def create_or_resume(**_):
    """Create and re-sync a Nova instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    for cell in CELLS:
        # NOTE(mnaser): cell0 does not need a message queue
        if cell != 'cell0':
            if not utils.ensure_secret("openstack", "nova-%s-rabbitmq" % cell):
                utils.create_or_update('nova/secret-rabbitmq.yml.j2',
                                       component=cell,
                                       password=utils.generate_password())
            utils.create_or_update('nova/rabbitmq.yml.j2', component=cell)
