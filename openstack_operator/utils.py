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
import json
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


BACKUP_SCHEDULE = {
    'magnum': '0 0 * * * *',
    'barbican': '0 5 * * * *',
    'cinder': '0 10 * * * *',
    'glance': '0 15 * * * *',
    'heat': '0 20 * * * *',
    'neutron': '0 25 * * * *',
    'octavia': '0 30 * * * *',
    'nova': '0 35 * * * *',
    'keystone': '0 40 * * * *',
    'zuul': '0 58 * * * *'
}

DIR_PATH = os.path.dirname(os.path.realpath(__file__))

VERSION = version.VersionInfo('openstack_operator').version_string()


def to_yaml(value):
    """Return a YAML string from a dictionary."""
    return yaml.safe_dump(value)


def to_dict(value):
    """Return a dictionary from a YAML string"""
    return yaml.safe_load(value)


def labels(app, instance=None, component=None):
    """Return standard labels for the operator."""
    metadata = {
        'app.kubernetes.io/managed-by': 'openstack-operator',
        'app.kubernetes.io/name': app,
    }
    if instance:
        metadata['app.kubernetes.io/instance'] = instance
    if component:
        metadata['app.kubernetes.io/component'] = component
    return yaml.safe_dump(metadata).strip()


ENV = jinja2.Environment(
    loader=jinja2.FileSystemLoader("%s/templates" % DIR_PATH)
)
ENV.filters['to_yaml'] = to_yaml
ENV.globals['labels'] = labels


def create_or_update(template, server_side=True, **kwargs):
    """Create or update a Kubernetes resource.

    This function is called with a template and the args to pass to that
    template and it will generate a Kubernetes object, with that
    object, it will try and check if it exists.  If it does, it will run
    an update, if not, it will create the new object.
    """

    resource = generate_object(template, **kwargs)

    if server_side:
        # NOTE(mnaser): The following relies on server-side apply and requires
        #               at least Kuberentes v1.16+
        resp = resource.api.patch(
            **resource.api_kwargs(
                headers={
                    "Content-Type": "application/apply-patch+yaml"
                },
                params={
                    'fieldManager': 'openstack-operator',
                    'force': True,
                },
                data=to_yaml(resource.obj),
            )
        )
        resource.api.raise_for_status(resp)
        resource.set_obj(resp.json())
    else:
        obj = copy.deepcopy(resource.obj)
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

    if "adopt" in kwargs and kwargs["adopt"]:
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

    try:
        secret = objects.Secret.objects(api).filter(namespace=namespace).get(
            name=name
        )
    except pykube.exceptions.ObjectDoesNotExist:
        return None

    return {
        k: base64.b64decode(v).decode('utf-8')
        for k, v in secret.obj['data'].items()
    }


def ensure_secret(namespace, name):
    """Check if a secret exists

    This function return true when the specific secret exists while
    return false when does not exist"""

    api = pykube.HTTPClient(pykube.KubeConfig.from_env())

    try:
        objects.Secret.objects(api).filter(namespace=namespace).get(
            name=name
        )
        return True
    except pykube.exceptions.ObjectDoesNotExist:
        return False


def generate_hash(dictionary):
    """Generate a hash from a dictionary, return None
    if dictionary is empty"""

    if not dictionary:
        return None
    return hash(json.dumps(dictionary))


def get_configmap(namespace, name):
    """Retrieve a configmap from Kubernetes.

    This function retrieves a configmap from Kubernetes, decodes it and passes
    the value of the data
    """

    api = pykube.HTTPClient(pykube.KubeConfig.from_env())

    try:
        config = objects.ConfigMap.objects(api).filter(
            namespace=namespace
        ).get(name=name)
    except pykube.exceptions.ObjectDoesNotExist:
        return None

    return config.obj["data"]


def get_backup_schedule(name):
    """Retrieve backup schedule for openstack services

    This function retrieves a backup schedule for the specified openstack
    service and the schedule is a cronjob format"""

    if name not in BACKUP_SCHEDULE:
        return "0 0 * * * *"
    return BACKUP_SCHEDULE[name]


def deploy_memcached(name, **_):
    """
    Deploy a generic instance of Memcached

    This function deploys a generic instance of Memcached with sane defaults,
    it's meant to be here to be consumed/called by the services.
    """
    create_or_update('operator/memcached.yml.j2', name=name)


def deploy_uwsgi_config():
    """Deploy a default configmap for uwsgi apps

    This function deploys a default configmap for uwsgi apps."""

    create_or_update('operator/uwsgidefaultconfig.yml.j2')


def deploy_rabbitmq(name, **_):
    """
    Deploy a generic instance of rabbitmq

    This function deploys a generic instance of Rabbitmq with a secret,
    it's meant to be here to be consumed/called by the services.
    The secret should include user and password.
    """

    if not ensure_secret("openstack", name + "-rabbitmq"):
        create_or_update('operator/secret-rabbitmq.yml.j2',
                         name=name, password=generate_password())
    create_or_update('operator/rabbitmq.yml.j2', name=name)


def ensure_mysql_cluster(name, spec=None, user=None, database=None):
    """Create or update mysql cluster"""

    if spec is None:
        spec = {}

    if database is None:
        database = name
    if user is None:
        user = database

    config = get_secret("openstack", name + "-mysql")
    if config is None:
        root_password = generate_password()
        password = generate_password()
        create_or_update('mysqlcluster/secret-mysqlcluster.yml.j2',
                         name=name, user=user,
                         database=database, password=password,
                         rootPassword=root_password)
        config = get_secret("openstack", name + "-mysql")

    config['connection'] = \
        "mysql+pymysql://%s:%s@%s-mysql-master/%s?charset=utf8" % (
            config["USER"],
            config["PASSWORD"],
            name,
            config["DATABASE"]
        )

    create_or_update('mysqlcluster/mysqlcluster.yml.j2',
                     server_side=False, name=name, spec=spec)
    return config
