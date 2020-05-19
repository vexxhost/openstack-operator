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
# pylint: disable=W0613
"""Openstack Operator

This module maintains the operator startup, it takes care of creating
the appropriate deployments, an instance of Keystone, Heat and Horizon
 for the installation.
"""

import os
import kopf

from openstack_operator import heat
from openstack_operator import horizon
from openstack_operator import keystone
from openstack_operator import utils


OPERATOR_CONFIGMAP = "operator-config"


def _create_namespace():
    """Create a namespace for the operator

    All resources which are managed by the operator would
    be deployed on this namespace"""

    utils.create_or_update('operator/namespace.yml.j2')


@kopf.on.startup()
async def startup_fn(logger, **kwargs):
    """Create several deployments at the startup of the operator

    keystone, heat, and horizon
    """

    namespace = os.getenv('OPERATOR_NAMESPACE')
    config = utils.get_configmap(namespace, OPERATOR_CONFIGMAP)
    config = utils.to_dict(config["operator-config.yaml"])
    _create_namespace()
    if "keystone" in config:
        keystone.create_or_resume("keystone", config["keystone"])
    if "horizon" in config:
        horizon.create_secret("horizon")
        horizon.create_or_resume("horizon", config["horizon"])
    if "heat" in config:
        heat.create_or_resume("heat", config["heat"])


@kopf.on.update('', 'v1', 'configmaps')
def update(name, namespace, new, **_):
    """Update the startup deployments when the operator configmap is changed

    keystone, heat, and horizon
    """
    if namespace == os.getenv('OPERATOR_NAMESPACE') \
       and name == OPERATOR_CONFIGMAP:
        config = utils.to_dict(new["data"]["operator-config.yaml"])
        if "horizon" in config:
            horizon.update("horizon", config["horizon"])
