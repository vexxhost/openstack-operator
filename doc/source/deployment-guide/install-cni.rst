Install CNI
-----------
The tested and supported CNI for the OpenStack operator is Calico due to it's
high performance and support for ``NetworkPolicy``.   You can deploy it onto
the cluster by running the following::

   $ iptables -I DOCKER-USER -j ACCEPT
   $ kubectl apply -f https://docs.opendev.org/vexxhost/openstack-operator/calico.yaml

.. note::

   The first command overrides Docker's behaviour of disabling all traffic
   routing if it is enabled, as this is necessary for the functioning on the
   Kubernetes cluster.

Once the CNI is deployed, you'll have to make sure that Calico detected the
correct interface to build your BGP mesh, you can run this command and make
sure that all systems are on the right network::

   $ kubectl describe nodes | grep IPv4Address

If they are not on the right IP range or interface, you can run the following
command and edit the ``calico-node`` DaemonSet::

   $ kubectl -n kube-system edit ds/calico-node

You'll need to add an environment variable to the container definition which
skips the interfaces you don't want, something similar to this::

   - name: IP_AUTODETECTION_METHOD
     value: skip-interface=bond0
