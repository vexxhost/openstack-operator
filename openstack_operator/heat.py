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

"""Heat Operator

This module maintains the operator for Heat, it takes care of creating
the appropriate deployments, an instance of Memcache, RabbitMQ and a database
server for the installation.
"""


from openstack_operator import identity
from openstack_operator import utils

MEMCACHED = True


def create_or_resume(name, spec, **_):
    """Create and re-sync any Heat instances
    """

    utils.ensure_mysql_cluster("heat", spec["mysql"])

    # deploy rabbitmq
    utils.deploy_rabbitmq("heat")

    config_hash = utils.generate_hash(spec)
    # deploy heat api
    utils.create_or_update('heat/api/daemonset.yml.j2', spec=spec,
                           config_hash=config_hash)
    utils.create_or_update('heat/api/service.yml.j2')

    # deploy heat cfn api
    utils.create_or_update('heat/api-cfn/daemonset.yml.j2', spec=spec,
                           config_hash=config_hash)
    utils.create_or_update('heat/api-cfn/service.yml.j2')

    # deploy heat cfn engine
    utils.create_or_update('heat/engine/daemonset.yml.j2', spec=spec,
                           config_hash=config_hash)

    # deploy clean jobs
    utils.create_or_update('heat/cronjob-service-clean.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('heat/cronjob-purge-deleted.yml.j2',
                           name=name, spec=spec)

    api_url = cfn_url = None
    if "ingress" in spec:
        utils.create_or_update('heat/ingress.yml.j2',
                               name=name, spec=spec)
        api_url = spec["ingress"]["host"]["api"]
        cfn_url = spec["ingress"]["host"]["api-cfn"]

    # Create application credential
    identity.ensure_application_credential(name="heat")

    # Create service and endpoints
    if "endpoint" not in spec:
        spec["endpoint"] = True
    if spec["endpoint"]:
        identity.ensure_service(name="heat-api",
                                service_type="orchestration",
                                url=api_url, path="/v1/$(project_id)s",
                                desc="Heat Orchestration Service")
        identity.ensure_service(name="heat-api-cfn",
                                service_type="cloudformation",
                                url=cfn_url, path="/v1",
                                desc="Heat CloudFormation Service")


def update(name, spec, **_):
    """Update a heat

    This function updates the deployment for heat if there are any
    changes that happen within it.
    """
    if "ingress" in spec:
        utils.create_or_update('horizon/ingress.yml.j2',
                               name=name, spec=spec)
