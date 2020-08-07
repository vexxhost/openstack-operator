#!/usr/local/bin/python
# Copyright (c) 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
WSGI config for openstack_dashboard project.
"""

import os
import sys

import pkg_resources
import sentry_sdk

from django.core.wsgi import get_wsgi_application
from sentry_sdk.integrations import wsgi

VERSION = pkg_resources.get_distribution("horizon").version

sentry_sdk.init(
    release="horizon@%s" % VERSION,
    traces_sample_rate=0.1
)

# Add this file path to sys.path in order to import settings
sys.path.insert(0, os.path.normpath(os.path.join(
    os.path.dirname(os.path.realpath(__file__)), '..')))
os.environ['DJANGO_SETTINGS_MODULE'] = 'openstack_dashboard.settings'
sys.stdout = sys.stderr

application = get_wsgi_application()
application = wsgi.SentryWsgiMiddleware(application)
