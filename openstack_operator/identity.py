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

import kopf

from openstack_operator import utils


def ensure_service(name, service, desc, url=None):
    """Create or update service and endpoints
    """

    try:
        # Create or resume service
        utils.create_or_update('identity/service.yml.j2', name=name,
                               type=service, description=desc)

        # Create or resume endpoints
        internal_url = "http://" + name + ".openstack.svc.cluster.local"
        public_url = internal_url
        if url is not None:
            public_url = "http://" + url
        utils.create_or_update('identity/endpoint.yml.j2',
                               service=service, interface='internal',
                               url=internal_url)
        utils.create_or_update('identity/endpoint.yml.j2',
                               service=service, interface='public',
                               url=public_url)
    except Exception as ex:
        raise kopf.TemporaryError(str(ex), delay=5)
