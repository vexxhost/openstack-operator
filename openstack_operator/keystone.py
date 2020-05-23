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

"""Keystone Operator

This module maintains the operator for Keystone which does everything from
deployment to taking care of rotating fernet & credentials keys."""

import base64
import kopf

from cryptography import fernet

from openstack_operator import filters
from openstack_operator import utils

TOKEN_EXPIRATION = 86400
FERNET_ROTATION_INTERVAL = 3600


def _is_keystone_deployment(name, **_):
    return name == 'keystone'


def create_or_rotate_fernet_repository(name):
    """Create or rotate fernet tokens

    This will happen when it sees a Keystone deployment that we manage and it
    will initialize (or rotate) the fernet repository.
    """

    data = utils.get_secret('openstack', 'keystone-%s' % (name))

    # Stage an initial key 0 if we don't have anything.
    if data is None:
        data = {'0': fernet.Fernet.generate_key().decode('utf-8')}

    # Get highest key number
    sorted_keys = [int(k) for k in data.keys()]
    sorted_keys.sort()
    next_key = str(max(sorted_keys) + 1)

    # Promote key 0 to primary
    data[next_key] = data['0']
    sorted_keys.append(int(next_key))

    # Stage a new key
    data['0'] = fernet.Fernet.generate_key().decode('utf-8')

    # Determine number of active keys
    active_keys = int(TOKEN_EXPIRATION / FERNET_ROTATION_INTERVAL)

    # Determine the keys to keep and drop others
    keys_to_keep = [0] + sorted_keys[-active_keys:]
    keys = {k: base64.b64encode(v.encode('utf-8')).decode('utf-8')
            for k, v in data.items() if int(k) in keys_to_keep}

    # Update secret
    utils.create_or_update('keystone/secret-fernet.yml.j2', name=name,
                           keys=keys, is_strategic=False, adopt=True)


@kopf.timer('apps', 'v1', 'deployments',
            when=kopf.all_([filters.managed, _is_keystone_deployment]),
            interval=FERNET_ROTATION_INTERVAL)
def create_or_rotate_fernet(**_):
    """Create or rotate fernet keys

    This will happen when it sees a Keystone deployment that we manage and it
    will initialize (or rotate) the fernet repository.
    """

    create_or_rotate_fernet_repository('fernet')
    create_or_rotate_fernet_repository('credential')


@kopf.on.resume('identity.openstack.org', 'v1alpha1', 'keystones')
@kopf.on.create('identity.openstack.org', 'v1alpha1', 'keystones')
def create_or_resume(name, spec, **_):
    """Create and re-sync any Keystone instances

    This function is called when a new resource is created but also when we
    start the service up for the first time.
    """
    env = utils.get_uwsgi_env()
    utils.create_or_update('keystone/deployment.yml.j2',
                           name=name, spec=spec, env=env)
    utils.create_or_update('keystone/service.yml.j2',
                           name=name, spec=spec)
    utils.create_or_update('keystone/horizontalpodautoscaler.yml.j2',
                           name=name)
