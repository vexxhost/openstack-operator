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

"""Kubernetes Objects

This module maintains a list of all of the Kubernetes objects that are used
by the operator.  It also includes a few of the custom ones that we use which
are not part of ``pykube``.

It also inclues a ``dict`` with mappings which allows doing reverse-lookups
from combinations of apiVersion and kind to the exact model.
"""

from pykube.objects import APIObject
from pykube.objects import ConfigMap
from pykube.objects import CronJob
from pykube.objects import DaemonSet
from pykube.objects import Deployment
from pykube.objects import HorizontalPodAutoscaler
from pykube.objects import Ingress
from pykube.objects import Namespace
from pykube.objects import NamespacedAPIObject
from pykube.objects import Pod
from pykube.objects import Secret
from pykube.objects import Service
from pykube.objects import StatefulSet


class IdentityApplicationCredential(APIObject):
    """ApplicationCredential Kubernetes object"""

    version = "identity.openstack.org/v1alpha1"
    endpoint = "applicationcredentials"
    kind = "ApplicationCredential"


class IdentityService(APIObject):
    """Service Kubernetes object"""

    version = "identity.openstack.org/v1alpha1"
    endpoint = "services"
    kind = "Service"


class IdentityEndpoint(APIObject):
    """Endpoint Kubernetes object"""

    version = "identity.openstack.org/v1alpha1"
    endpoint = "endpoints"
    kind = "Endpoint"


class Mcrouter(NamespacedAPIObject):
    """Mcrouter Kubernetes object"""

    version = "infrastructure.vexxhost.cloud/v1alpha1"
    endpoint = "mcrouters"
    kind = "Mcrouter"


class Memcached(NamespacedAPIObject):
    """Memcached Kubernetes object"""

    version = "infrastructure.vexxhost.cloud/v1alpha1"
    endpoint = "memcacheds"
    kind = "Memcached"


class MysqlCluster(NamespacedAPIObject):
    """Mysql Cluster Kubernetes object"""

    version = "mysql.presslabs.org/v1alpha1"
    endpoint = "mysqlclusters"
    kind = "MysqlCluster"


class PodMonitor(NamespacedAPIObject):
    """PodMonitor Kubernetes object"""

    version = "monitoring.coreos.com/v1"
    endpoint = "podmonitors"
    kind = "PodMonitor"


class PrometheusRule(NamespacedAPIObject):
    """PrometheusRule Kubernetes object"""

    version = "monitoring.coreos.com/v1"
    endpoint = "prometheusrules"
    kind = "PrometheusRule"


class Rabbitmq(NamespacedAPIObject):
    """Rabbitmq Kubernetes object"""

    version = "infrastructure.vexxhost.cloud/v1alpha1"
    endpoint = "rabbitmqs"
    kind = "Rabbitmq"


MAPPING = {
    "v1": {
        "ConfigMap": ConfigMap,
        "Namespace": Namespace,
        "Pod": Pod,
        "Secret": Secret,
        "Service": Service,
    },
    "apps/v1": {
        "DaemonSet": DaemonSet,
        "Deployment": Deployment,
        "StatefulSet": StatefulSet,
    },
    "autoscaling/v1": {
        "HorizontalPodAutoscaler": HorizontalPodAutoscaler,
    },
    "batch/v1beta1": {
        "CronJob": CronJob,
    },
    "extensions/v1beta1": {
        "Ingress": Ingress
    },
    "identity.openstack.org/v1alpha1": {
        "ApplicationCredential": IdentityApplicationCredential,
        "Service": IdentityService,
        "Endpoint": IdentityEndpoint
    },
    "infrastructure.vexxhost.cloud/v1alpha1": {
        "Mcrouter": Mcrouter,
        "Memcached": Memcached,
        "Rabbitmq": Rabbitmq
    },
    "monitoring.coreos.com/v1": {
        "PodMonitor": PodMonitor,
        "PrometheusRule": PrometheusRule,
    },
    "networking.k8s.io/v1beta1": {
        "Ingress": Ingress
    },
    "mysql.presslabs.org/v1alpha1": {
        "MysqlCluster": MysqlCluster
    }
}
