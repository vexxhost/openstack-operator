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

"""Endpoints Operator

This operator helps manage the creation and removal of endpoints inside
Keystone using custom resources.
"""

import kopf
from openstack_operator import utils


def _get_service_by_type(conn, service_type):
    """Get a service from Keystone based on service type."""

    services = conn.search_services(filters={"type": service_type})

    if len(services) > 1:
        raise RuntimeError("Multiple services with type: %s" % service_type)
    if len(services) == 0:
        raise RuntimeError("Unable to find service: %s" % service_type)
    return services[0]


def _get_endpoint(conn, service_type, interface):
    """Get an endpoint from Keystone

    This method will retrieve the endpoint from Keystone, raise an error if it
    found more than one or return None if it couldn't find it
    """

    service = _get_service_by_type(conn, service_type)

    filters = {
        "service_id": service.id,
        "interface": interface,
        "region": conn.config.get_region_name(),
    }
    endpoints = conn.search_endpoints(filters=filters)

    if len(endpoints) > 1:
        raise RuntimeError("Found multiple endpoints with interface & region")
    if len(endpoints) == 0:
        return service, None
    return service, endpoints[0]


@kopf.on.resume('identity.openstack.org', 'v1alpha1', 'endpoints')
@kopf.on.create('identity.openstack.org', 'v1alpha1', 'endpoints')
def create_or_resume(spec, **_):
    """Create or resume controller

    This function runs when a new resource is created or when the controller
    is first started.  It creates or updates the appropriate endpoint.
    """

    conn = utils.get_openstack_connection()
    service, endpoint = _get_endpoint(conn, spec["service"], spec["interface"])
    if endpoint:
        conn.update_endpoint(endpoint.id, url=spec["url"])
        return

    conn.create_endpoint(service_name_or_id=service.id, url=spec["url"],
                         interface=spec["interface"],
                         region=conn.config.get_region_name())


@kopf.on.delete('identity.openstack.org', 'v1alpha1', 'endpoints')
def delete(spec, **_):
    """Delete an endpoint

    This function runs when the endpoint CR is deleted and removes the record
    from Keystone.
    """

    conn = utils.get_openstack_connection()
    endpoint = _get_endpoint(conn, spec["service"], spec["interface"])

    if not endpoint:
        return

    conn.delete_endpoint(endpoint)
