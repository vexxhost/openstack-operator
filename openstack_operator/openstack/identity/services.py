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

"""Services Operator

This operator helps manage the creation and removal of services inside
Keystone using custom resources.
"""

import kopf

from openstack_operator import utils


def _get_service(conn, name, service_type):
    """Get a service from Keystone

    This method will retrieve the service from Keystone, raise an error if it
    found more than one or return None if it couldn't find it
    """

    services = conn.search_services(name_or_id=name,
                                    filters={"type": service_type})

    if len(services) > 1:
        raise RuntimeError("Found multiple services with name and type")
    if len(services) == 0:
        return None
    return services[0]


@kopf.on.resume('identity.openstack.org', 'v1alpha1', 'services')
@kopf.on.create('identity.openstack.org', 'v1alpha1', 'services')
def create_or_resume(name, spec, **_):
    """Create or resume controller

    This function runs when a new resource is created or when the controller
    is first started.  It creates or updates the appropriate service.
    """

    conn = utils.get_openstack_connection()
    service = _get_service(conn, name, spec["type"])

    if service:
        service = conn.update_service(service.id, name=name,
                                      type=spec["type"],
                                      description=spec["description"])
        return

    service = conn.create_service(name=name, type=spec["type"],
                                  description=spec["description"])


@kopf.on.delete('identity.openstack.org', 'v1alpha1', 'services')
def delete(name, spec, **_):
    """Delete a service

    This function runs when the servce CR is deleted and removes the record
    from Keystone.
    """

    conn = utils.get_openstack_connection()
    service = _get_service(conn, name, spec["type"])

    if not service:
        return

    conn.delete_service(service)
