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

"""glance Operator

This module maintains the operator for Mcrouter, it takes care of creating
the appropriate deployments, Mcrouter, pod monitors and Prometheus rules.
"""


from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync a glance instance

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    # deploy mysql for glance
    utils.ensure_mysql_cluster("glance", spec["mysql"])

    # deploy glance api
    utils.create_or_update('glance/daemonset.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('glance/service.yml.j2',
                           name=name, spec=spec)
    url = None
    if "ingress" in spec:
        utils.create_or_update('glance/ingress.yml.j2',
                               name=name, spec=spec)
        url = spec["ingress"]["host"]

    # Create application credential
    identity.ensure_application_credential(name="glance")

    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="glance", service_type="image",
                                url=url, desc="Glance Image Service")


def update(name, spec, **_):
    """Update a glance

    This function updates the deployment for glance if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('glance/ingress.yml.j2',
                               name=name, spec=spec)
