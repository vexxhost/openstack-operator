Setup Load Balancer
-------------------
The load balancer which will be distributing requests across all of the
Kubernetes API servers will be HAproxy.

.. note::

   We do not suggest using HAproxy to distribute load across all of the
   ingress controllers.  The primary reason being that it introduces an extra
   hop in the network for no large benefit.  The ingress should be bound
   directly on the virtual IP.

The following example assumes that you have 3 controllers, with their IP
addresses being ``10.0.0.1``, ``10.0.0.2``, ``10.0.0.3``.  It also assumes that
all of the Kubernetes API servers will be listening on port ``16443`` and it
will be listening on port ``6443``.

You'll have to create a configuration file on the local system first::

    $ mkdir /etc/haproxy
    $ cat <<EOF | tee /etc/haproxy/haproxy.cfg
    listen kubernetes
      mode tcp
      bind 0.0.0.0:6443
      timeout connect 30s
      timeout client 4h
      timeout server 4h
      server ctl1 10.0.0.1:16443 check
      server ctl2 10.0.0.2:16443 check
      server ctl3 10.0.0.3:16443 check
    EOF

Once you've setup the configuration file, you can start up the containerized
instance of HAproxy::

   $ docker run --net=host \
                --volume=/etc/haproxy:/usr/local/etc/haproxy:ro \
                --detach \
                --restart always \
                --name=haproxy \
                haproxy:2.2

You'll also need to make sure that you have a DNS record pointing towards your
virtual IP address.  It is also recommended that you create a wildcard DNS as
well to allow multiple hosts for the ingress without needing extra changes in
your DNS, something like this::

   cloud.vexxhost.net.  86400   IN	A	10.0.0.200
   *.cloud.vexxhost.net 86400   IN      A       10.0.0.200
