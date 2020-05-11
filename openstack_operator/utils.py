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

"""Utilities

The module contains a few useful utilities which we refactor out in order
to be able to use them across all different operators.
"""
import base64
import copy
import operator
import os
import secrets
import string

import jinja2
import kopf
from pbr import version
import pykube
import yaml
import openstack

from openstack_operator import objects


DIR_PATH = os.path.dirname(os.path.realpath(__file__))

UWSGI_SETTINGS = {
    'UWSGI_ENABLE_THREADS': True,
    'UWSGI_PROCESSES': 2,
    'UWSGI_EXIT_ON_RELOAD': True,
    'UWSGI_DIE_ON_TERM': True,
    'UWSGI_LAZY_APPS': True,
    'UWSGI_ADD_HEADER': 'Connection: close',
    'UWSGI_BUFFER_SIZE': 65535,
    'UWSGI_THUNDER_LOCK': True,
    'UWSGI_DISABLE_LOGGING': True,
    'UWSGI_AUTO_CHUNCKED': True,
    'UWSGI_HTTP_RAW_BODY': True,
    'UWSGI_SOCKET_TIMEOUT': 10,
    'UWSGI_NEED_APP': True,
    'UWSGI_ROUTE_USER_AGENT': '^kube-probe.* donotlog:'
}

VERSION = version.VersionInfo('openstack_operator').version_string()


def to_yaml(value):
    """Return a YAML string from a dictionary."""
    return yaml.safe_dump(value)


def to_dict(value):
    """Return a dictionary from a YAML string"""
    return yaml.safe_load(value)


def labels(app, instance, component=None):
    """Return standard labels for the operator."""
    metadata = {
        'app.kubernetes.io/managed-by': 'openstack-operator',
        'app.kubernetes.io/name': app,
        'app.kubernetes.io/instance': instance,
    }
    if component:
        metadata['app.kubernetes.io/component'] = component
    return yaml.safe_dump(metadata).strip()


ENV = jinja2.Environment(
    loader=jinja2.FileSystemLoader("%s/templates" % DIR_PATH)
)
ENV.filters['to_yaml'] = to_yaml
ENV.globals['labels'] = labels


def create_or_update(template, **kwargs):
    """Create or update a Kubernetes resource.

    This function is called with a template and the args to pass to that
    template and it will generate a Kubernetes object, with that
    object, it will try and check if it exists.  If it does, it will run
    an update, if not, it will create the new object.
    """

    resource = generate_object(template, **kwargs)
    obj = copy.deepcopy(resource.obj)

    # Try to get the remote record
    try:
        resource.reload()
        resource.obj = obj
        resource.update()
    except pykube.exceptions.HTTPError as exc:
        if exc.code != 404:
            raise
        resource.create()

    return resource


def ensure_absent(template, **kwargs):
    """Ensure a Kubernetes resource bound to a template is deleted

    This function gets a template and makes sure that the object doesn't
    exist on the remote cluster.
    """

    resource = generate_object(template, **kwargs)
    resource.delete()


def render_template(template, **kwargs):
    """Render template from YAML files.

    This function renders a template based on provided keyword arguments.
    """

    template = ENV.get_template(template)
    yamldoc = template.render(**kwargs)
    return yaml.safe_load(yamldoc)


def generate_yaml(template, **kwargs):
    """Generate dictionary from YAML template.

    This takes a Jinja2 template, renders it using all the passed ``kwargs``
    and then runs ``adopt`` as well to prepare it to be committed to the
    cluster.
    """

    doc = render_template(template, **kwargs)
    kopf.adopt(doc)

    return doc


def generate_object(template, **kwargs):
    """Generate Kubernetes object

    This function renders a Jinja2 template provided into a Kubernetes object
    based on the ``apiVersion`` and ``kind`` from the generated YAML file.
    """

    doc = generate_yaml(template, **kwargs)

    api_version = doc['apiVersion']
    kind = doc['kind']
    resource = objects.MAPPING[api_version][kind]

    api = pykube.HTTPClient(pykube.KubeConfig.from_env())
    return resource(api, doc)


def get_ready_pod_ips(namespace, selector):
    """Get list of all ready pod IPs.

    This is a helper function which given a selector, will retrieve all pods
    and return the IP addresses of the ones that are ready.
    """

    api = pykube.HTTPClient(pykube.KubeConfig.from_env())
    pods = objects.Pod.objects(api).filter(namespace=namespace,
                                           selector=selector)
    ready_pods = filter(operator.attrgetter("ready"), pods)
    servers = sorted([p.obj["status"]["podIP"] for p in ready_pods])

    return servers


def get_openstack_connection():
    """Get an instance of OpenStack SDK."""
    return openstack.connect(cloud="envvars", app_name='openstack-operator',
                             app_version=VERSION)


def generate_password(length=20):
    """Generate a random password."""

    alphabet = string.ascii_letters + string.digits
    return ''.join(secrets.choice(alphabet) for i in range(length))


def get_secret(namespace, name):
    """Retrieve a secret from Kubernetes.

    This function retrieves a Secret from Kubernetes, decodes it and passes
    the value of the data
    """

    api = pykube.HTTPClient(pykube.KubeConfig.from_env())

    secret = objects.Secret.objects(api).filter(namespace=namespace).get(
        name=name
    )

    return {
        k: base64.b64decode(v).decode('utf-8')
        for k, v in secret.obj['data'].items()
    }


def generate_hash(dictionary):
    """Generate a hash from a dictionary, return None if dictionary is empty"""

    if not dictionary:
        return None
    return hash(frozenset(dictionary.items()))


def get_uwsgi_env():
    """Generate k8s env list from UWSGI_SETTINGS dict"""
    res = []
    for key, value in UWSGI_SETTINGS.items():
        res.append({'name': key, 'value': value})
    return res
