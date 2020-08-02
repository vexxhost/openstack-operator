Setup Virtual IP
----------------
The virtual IP runs across all controllers in order to allow the Kubernetes API
server to be highly available and load balancer.  It also becomes one of the
interfaces that the Kubernetes ingress listens on where all the OpenStack API
endpoints will be exposed.

The recommended way of deploying a virtual IP address is using ``keepalived``
running inside Docker in order to make sure that your environment remains clean
and easily reproducible.

You should use the following command in order to start up ``keepalived`` to
host the virtual IP address.  These commands should be ran on all your
controllers and they assume that you have 3 controllers with IP addresses
``10.0.0.1``, ``10.0.0.2``, ``10.0.0.3``.  The following example is what you
would run on the ``10.0.0.1`` machine with a VIP of ``10.0.0.200`` running
on the interface ``eth0``::

   $ docker run --cap-add=NET_ADMIN \
                --cap-add=NET_BROADCAST \
                --cap-add=NET_RAW \
                --net=host \
                --env KEEPALIVED_INTERFACE=eth0 \
                --env KEEPALIVED_UNICAST_PEERS="#PYTHON2BASH:['10.0.0.2', '10.0.0.3']" \
                --env KEEPALIVED_VIRTUAL_IPS="#PYTHON2BASH:['10.0.0.200']" \
                --detach \
                --restart always \
                --name keepalived \
                osixia/keepalived:2.0.20

.. note::

   You'll have to make sure to edit the ``KEEPALIVED_UNICAST_PEERS``
   environment variable accordingly depending on the host you're running this
   on.  It should always point at the other hosts.


