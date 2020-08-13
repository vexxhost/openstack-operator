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

"""identity Operator

This module contains a few common functions for identity management
"""

from openstack_operator import utils


def ensure_service(name, service_type, desc, url=None, path=""):
    """Create or update service and endpoints

    name: service name
    service_type: service type
    desc: service descriptioin
    url: hostname of public endpoint
    path: sub path of endpoint
    """

    # Create or resume service
    utils.create_or_update('identity/service.yml.j2', name=name,
                           type=service_type, description=desc)

    # Create or resume endpoints
    internal_url = public_url = \
        "http://" + name + ".openstack.svc.cluster.local" + path

    if url is not None:
        public_url = "https://" + url + path
    utils.create_or_update('identity/endpoint.yml.j2',
                           service=service_type, interface='internal',
                           url=internal_url)
    utils.create_or_update('identity/endpoint.yml.j2',
                           service=service_type, interface='public',
                           url=public_url)


def ensure_application_credential(name):
    """Create or update applicationcredentials
    """

    utils.create_or_update('identity/applicationcredential.yml.j2',
                           name=name)
