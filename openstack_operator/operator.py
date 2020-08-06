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
import sentry_sdk
from sentry_sdk.integrations import aiohttp

from openstack_operator import ceilometer
from openstack_operator import chronyd
from openstack_operator import glance
from openstack_operator import heat
from openstack_operator import horizon
from openstack_operator import keystone
from openstack_operator import libvirtd_exporter
from openstack_operator import magnum
from openstack_operator import utils


OPERATOR_CONFIGMAP = "operator-config"

sentry_sdk.init(
    integrations=[aiohttp.AioHttpIntegration()],
    traces_sample_rate=1.0
)


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
        spec = set_service_config(config, "keystone")
        keystone.create_or_resume("keystone", spec)
    if "horizon" in config:
        spec = set_service_config(config, "horizon")
        horizon.create_or_resume("horizon", spec)
    if "heat" in config:
        spec = set_service_config(config, "heat")
        heat.create_or_resume("heat", spec)
    if "glance" in config:
        spec = set_service_config(config, "glance")
        glance.create_or_resume("glance", spec)
    if "magnum" in config:
        spec = set_service_config(config, "magnum")
        magnum.create_or_resume("magnum", spec)
    if "ceilometer" in config:
        spec = config["ceilometer"]
        ceilometer.create_or_resume(spec)

    if "chronyd" in config:
        spec = config["chronyd"]
    else:
        spec = {}
    chronyd.create_or_resume(spec)

    if "libvirtd-exporter" in config:
        spec = config["libvirtd-exporter"]
    else:
        spec = {}
    libvirtd_exporter.create_or_resume(spec)


def set_service_config(all_config, service_name):
    """Retrieve the config for each openstack service

    The config for each service is comprised of service-level
    config and operator-level config"""

    # Set the service level config
    spec = all_config[service_name]

    # Inject the operator level config to service level
    # Backup config for mysql
    all_config["backup"]["schedule"] = utils.get_backup_schedule(service_name)
    if "mysql" in spec:
        spec["mysql"].update(all_config["backup"])
    else:
        spec["mysql"] = all_config["backup"]

    return spec
