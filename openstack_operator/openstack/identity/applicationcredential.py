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

"""Application Credential Operator

This operator helps manage the creation and removal of application
credential inside Keystone using custom resources.
"""

import kopf
from openstack_operator import utils


def _get_admin_user_id():
    """Get admin user id"""

    conn = utils.get_openstack_connection()
    user_name = conn.config.auth["username"]
    domain_id = conn.config.auth["user_domain_id"]
    user = conn.get_user(name_or_id=user_name, domain_id=domain_id)
    return user.id


@kopf.on.resume('identity.openstack.org', 'v1alpha1', 'applicationcredentials')
@kopf.on.create('identity.openstack.org', 'v1alpha1', 'applicationcredentials')
def create_or_resume(name, **_):
    """Create or resume controller

    This function runs when a new resource is created or when the
    controller is first started.  It creates or updates the appropriate
    applicationcredential."""

    identity = utils.get_openstack_connection().identity

    user = _get_admin_user_id()
    credential = \
        identity.find_application_credential(user=user, name_or_id=name)

    if credential is None:
        credential = \
            identity.create_application_credential(user=user, name=name)
        utils.create_or_update(
            'identity/secret-applicationcredential.yml.j2',
            name=name, secret=credential.secret,
            id=credential.id, adopt=True)
        return

    # NOTE(Alex): Sometimes, double POST application_credential requests
    # are made to keystone API at the "same time".
    # The credential secret is not created in this case.
    # The following codes should fix this case.
    if utils.get_secret(name=name+"-application-credential",
                        namespace="openstack") is None:
        utils.create_or_update(
            'identity/secret-applicationcredential.yml.j2',
            name=name, secret=credential.secret,
            id=credential.id, adopt=True)


@kopf.on.delete('identity.openstack.org', 'v1alpha1', 'applicationcredentials')
def delete(name, **_):
    """Delete an endpoint

    This function runs when the applicationcredential CR is deleted and
    removes the record from Keystone.
    """

    identity = utils.get_openstack_connection().identity

    user = _get_admin_user_id()
    credential = \
        identity.find_application_credential(user=user, name_or_id=name)

    if credential is None:
        return

    identity.delete_application_credential(user=user,
                                           application_credential=name)
