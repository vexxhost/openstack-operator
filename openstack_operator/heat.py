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

import kopf

from openstack_operator import utils


@kopf.on.resume('orchestration.openstack.org', 'v1alpha1', 'heats')
@kopf.on.create('orchestration.openstack.org', 'v1alpha1', 'heats')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Heat instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """

    env = utils.get_uwsgi_env()
    for component in ("api", "api-cfn"):
        utils.create_or_update('heat/deployment.yml.j2',
                               name=name, spec=spec,
                               component=component, env=env)
        utils.create_or_update('heat/service.yml.j2',
                               name=name, component=component)

    utils.create_or_update('heat/deployment.yml.j2',
                           name=name, spec=spec, component='engine')
