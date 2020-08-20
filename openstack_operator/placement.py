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
Placement service

This code takes care of doing the operations of the OpenStack Placement API
service.
"""

from openstack_operator import database
from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync a placement instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    # deploy mysql for placement
    database.ensure_mysql_cluster("placement", spec=spec["mysql"])

    # deploy placement api
    utils.create_or_update('placement/daemonset.yml.j2', spec=spec)
    utils.create_or_update('placement/service.yml.j2', spec=spec)

    # Create application credential
    identity.ensure_application_credential(name="placement")

    url = None
    if "ingress" in spec:
        utils.create_or_update('placement/ingress.yml.j2',
                               name=name, spec=spec)
        url = spec["ingress"]["host"]

    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="placement", service_type="placement",
                                url=url, desc="Placement Service")


def update(name, spec, **_):
    """Update a placement

    This function updates the deployment for placement if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('placement/ingress.yml.j2',
                               name=name, spec=spec)
