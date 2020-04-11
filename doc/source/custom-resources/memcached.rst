.. _memcached:

Memcached
#########

Memcached is an in-memory key-value store for small chunks of arbitrary
data (strings, objects) from results of database calls, API calls, or page
rendering.

Architecture
************

This resource creates a ``Deployment`` with a hard coded replica count of two,
the size of every replica corresponds to half the size provided inside the
custom resource.

This resource does not expose a headless service, instead, it creates a managed
resource of :ref:`Mcrouter` which is automatically updated with the IPs of the
pods that are running ``memcached``.

Usage
*****

.. literalinclude :: ../../../config/samples/infrastructure_v1alpha1_memcached.yaml
   :language: yaml
