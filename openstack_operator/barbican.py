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

"""barbican Operator

This module maintains the operator for Mcrouter, it takes care of creating
the appropriate deployments, Mcrouter, pod monitors and Prometheus rules.
"""


from openstack_operator import database
from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync a barbican instance
    """

    # deploy mysql for barbican
    if "mysql" not in spec:
        database.ensure_mysql_cluster("barbican", {})
    else:
        database.ensure_mysql_cluster("barbican", spec["mysql"])

    # deploy barbican api
    utils.create_or_update('barbican/daemonset.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('barbican/service.yml.j2',
                           name=name, spec=spec)

    url = None
    if "ingress" in spec:
        utils.create_or_update('barbican/ingress.yml.j2',
                               name=name, spec=spec)
        url = spec["ingress"]["host"]
    identity.ensure_service(name="barbican", service_type="key-manager",
                            url=url, desc="Barbican Service")


def update(name, spec, **_):
    """Update a barbican

    This function updates the deployment for barbican if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('barbican/ingress.yml.j2',
                               name=name, spec=spec)
