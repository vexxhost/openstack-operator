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

from openstack_operator import ceilometer
from openstack_operator import chronyd
from openstack_operator import heat
from openstack_operator import horizon
from openstack_operator import keystone
from openstack_operator import magnum
from openstack_operator import utils


OPERATOR_CONFIGMAP = "operator-config"


def operator_configmap(namespace, name, **_):
    """Filter on the operator's ConfigMap."""

    return namespace == os.getenv('OPERATOR_NAMESPACE', 'default') \
        and name == OPERATOR_CONFIGMAP


@kopf.on.resume('', 'v1', 'configmaps', when=operator_configmap)
@kopf.on.create('', 'v1', 'configmaps', when=operator_configmap)
@kopf.on.update('', 'v1', 'configmaps', when=operator_configmap)
def deploy(name, namespace, new, **_):
    """Update the startup deployments when the operator configmap is changed

    keystone, heat, and horizon
    """

    utils.create_or_update('operator/namespace.yml.j2')
    utils.create_or_update('operator/uwsgidefaultconfig.yml.j2')

    config = utils.to_dict(new['data']['operator-config.yaml'])

    if "keystone" in config:
        keystone.create_or_resume("keystone", config["keystone"])
    if "horizon" in config:
        horizon.create_secret("horizon")
        horizon.create_or_resume("horizon", config["horizon"])
    if "heat" in config:
        heat.create_or_resume("heat", config["heat"])
    if "magnum" in config:
        magnum.create_or_resume("magnum", config["magnum"])
    if "chronyd" in config:
        chronyd.create_or_resume(config["chronyd"])
    if "ceilometer" in config:
        ceilometer.create_or_resume(config["ceilometer"])
