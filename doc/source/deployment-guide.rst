Deployment Guide
================

The OpenStack operator requires that you have a functional Kuberentes cluster
in order to be able to deploy it.  The following steps out-line the
installation of a cluster and the operator.

The deployment of the OpenStack operator is highly containerised, even for the
components not managed by the operator.  The steps to get started involve
deploying Docker for running the underlying infrastructure, installing a
load balancer to access the Kubernetes API, deploying Kubernetes itself and
then the operator which should start the OpenStack deployment.

.. highlight:: console
.. include:: deployment-guide/install-docker.rst
.. include:: deployment-guide/setup-virtual-ip.rst
.. include:: deployment-guide/setup-load-balancer.rst
.. include:: deployment-guide/install-kubernetes.rst
.. include:: deployment-guide/install-cni.rst
