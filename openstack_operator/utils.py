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

import copy
import operator
import os

import jinja2
import kopf
import pykube
import yaml

from openstack_operator import objects


DIR_PATH = os.path.dirname(os.path.realpath(__file__))


def to_yaml(value):
    """Return a YAML string from a dictionary."""
    return yaml.safe_dump(value)


def labels(app, instance):
    """Return standard labels for the operator."""
    return yaml.safe_dump({
        'app.kubernetes.io/managed-by': 'openstack-operator',
        'app.kubernetes.io/name': app,
        'app.kubernetes.io/instance': instance,
    }).strip()


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


def generate_yaml(template, **kwargs):
    """Generate dictionary from YAML template.

    This takes a Jinja2 template, renders it using all the passed ``kwargs``
    and then runs ``adopt`` as well to prepare it to be committed to the
    cluster.
    """

    template = ENV.get_template(template)
    yamldoc = template.render(**kwargs)
    doc = yaml.safe_load(yamldoc)
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
