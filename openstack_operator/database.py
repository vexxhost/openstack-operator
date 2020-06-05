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

"""database Operator

This module contains a few common functions for database management
"""

from openstack_operator import utils


def ensure_mysql_cluster(name, spec):
    """Create or update mysql cluster"""

    config = utils.get_secret("openstack", name + "-mysql")
    if config is None:
        root_password = utils.generate_password()
        password = utils.generate_password()
        user = name
        database = name
        utils.create_or_update('mysqlcluster/secret-mysqlcluster.yml.j2',
                               name=name, user=user,
                               database=database, password=password,
                               rootPassword=root_password)
        config = utils.get_secret("openstack", name + "-mysql")

    utils.create_or_update('mysqlcluster/mysqlcluster.yml.j2',
                           server_side=False, name=name, spec=spec)
    return config
