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
import pkg_resources

import kopf
import sentry_sdk
from sentry_sdk.integrations import aiohttp

from openstack_operator import barbican
from openstack_operator import ceilometer
from openstack_operator import chronyd
from openstack_operator import cinder
from openstack_operator import glance
from openstack_operator import heat
from openstack_operator import horizon
from openstack_operator import keystone
from openstack_operator import libvirtd_exporter
from openstack_operator import magnum
from openstack_operator import neutron
from openstack_operator import placement
from openstack_operator import utils


OPERATOR_CONFIGMAP = "operator-config"
VERSION = pkg_resources.get_distribution("openstack_operator").version

sentry_sdk.init(
    release="openstack-operator@%s" % VERSION,
    integrations=[aiohttp.AioHttpIntegration()],
    traces_sample_rate=1.0
)


def operator_configmap(namespace, name, **_):
    """Filter on the operator's ConfigMap."""

    return namespace == os.getenv('OPERATOR_NAMESPACE', 'default') \
        and name == "operator-config"


@kopf.on.event('', 'v1', 'configmaps', when=operator_configmap)
async def deploy_memcacheds(body, **_):
    """
    Deploy multiple Memcached instances for OpenStack services

    This function makes sure that Memcached is deployed for all services which
    use it when when the operator sees any changes to the configuration.
    """
    services = utils.to_dict(body['data']['operator-config.yaml']).keys()

    for entry_point in pkg_resources.iter_entry_points('operators'):
        if entry_point.name not in services:
            continue

        module = entry_point.load()
        if hasattr(module, 'MEMCACHED') and module.MEMCACHED:
            utils.deploy_memcached(entry_point.name)


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
    if "placement" in config:
        spec = set_service_config(config, "placement")
        neutron.create_or_resume(spec)
    if "neutron" in config:
        spec = set_service_config(config, "neutron")
        placement.create_or_resume("neutron", spec)
    if "horizon" in config:
        spec = set_service_config(config, "horizon")
        horizon.create_or_resume("horizon", spec)
    if "heat" in config:
        spec = set_service_config(config, "heat")
        heat.create_or_resume("heat", spec)
    if "glance" in config:
        spec = set_service_config(config, "glance")
        glance.create_or_resume("glance", spec)
    if "cinder" in config:
        spec = set_service_config(config, "cinder")
        cinder.create_or_resume("cinder", spec)
    if "magnum" in config:
        spec = set_service_config(config, "magnum")
        magnum.create_or_resume("magnum", spec)
    if "barbican" in config:
        spec = config["barbican"]
        barbican.create_or_resume("barbican", spec)
    if "ceilometer" in config:
        spec = config["ceilometer"]
        ceilometer.create_or_resume(spec)

    spec = config.get("chronyd", {})
    chronyd.create_or_resume(spec)

    spec = config.get("libvirtd_exporter", {})
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
