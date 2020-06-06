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

"""Ceilometer

The following modules maintains all the operations for deploying Ceilometer
alongside Atmosphere.
"""

from openstack_operator import utils


def create_or_resume(spec):
    """Create or start-up Ceilometer."""

    config_hash = utils.generate_hash(spec)

    utils.create_or_update('ceilometer/secret.yml.j2', spec=spec)
    utils.create_or_update('ceilometer/deployment-agent-notification.yml.j2',
                           spec=spec, config_hash=config_hash)
