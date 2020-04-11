.. _mcrouter:

Mcrouter
########

Mcrouter is a memcached protocol router for scaling Memcached deployments. It
is a core component of cache infrastructure at Facebook and Instagram where
Mcrouter handles almost 5 billion requests per second at peak.

Architecture
************

This resource creates a ``Deployment`` with a hard coded replica count of two
which hosts two identical instances of ``mcrouter`` for redundancy purposes.

It also exposes a ``Service`` resource which points at the two ``mcrouter``
instances running.

Usage
*****

.. literalinclude :: ../../../config/samples/infrastructure_v1alpha1_mcrouter.yaml
   :language: yaml
