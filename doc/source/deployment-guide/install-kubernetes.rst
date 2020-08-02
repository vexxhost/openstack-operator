Install Kubernetes
------------------
The recommended container runtime for the operator is ``containerd``, it is
also what is used in production.  This document outlines the installation of
Kubernetes using ``kubeadm``.  You'll need to start by installing the
Kubernetes components on all of the systems.

.. tabs::

   .. code-tab:: console Debian

      $ curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
      $ sudo add-apt-repository "deb https://apt.kubernetes.io/ kubernetes-xenial main"
      $ apt-get update
      $ apt-get install -y kubelet kubeadm kubectl
      $ apt-mark hold containerd.io kubelet kubeadm kubectl
      $ containerd config default > /etc/containerd/config.toml
      $ systemctl restart containerd

Once this is done, you'll need to start off by preparing the configuration file
for ``kubeadm``, which should look somethig like this::

   $ cat <<EOF | tee /etc/kubernetes/kubeadm.conf
   ---
   apiVersion: kubeadm.k8s.io/v1beta2
   kind: InitConfiguration
   localAPIEndpoint:
     bindPort: 16443
   nodeRegistration:
     criSocket: /run/containerd/containerd.sock
   ---
   apiVersion: kubeadm.k8s.io/v1beta2
   kind: ClusterConfiguration
   controlPlaneEndpoint: "cloud.vexxhost.net:6443"
   apiServer:
     extraArgs:
       oidc-issuer-url: https://accounts.google.com
       oidc-username-claim: email
       oidc-client-id: 1075333334902-iqnm5nbme0c36eir9gub5m62e6pbkbqe.apps.googleusercontent.com
   networking:
     podSubnet: 10.244.0.0/16
   EOF

.. note::

   The ``cloud.vexxhost.net`` should be replaced by the DNS address that you
   created in the previous step.

   The options inside ``extraArgs`` are there to allow for OIDC authentication
   via Google Suite.  You should remove those or replace them with your own
   OIDC provider.

   The pod subnet listed there is the one recommended for usage with Calico,
   which is the supported and tested CNI.

At this point, you should be ready to start and bring up your first control
plane node, you can execute the following on any of the controllers::

   $ kubeadm init --config /etc/kubernetes/kubeadm.conf --upload-certs

At that point, the cluster will be up and it's best to add the
``cluster-admin`` credentials into the ``root`` user for future management::

   $ mkdir -p $HOME/.kube
   $ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
   $ sudo chown $(id -u):$(id -g) $HOME/.kube/config

.. warning::

   For all of the commands running on other nodes, you'll need to make sure
   that you include the following flag or it will use Docker instead of the
   recommended ``containerd``::

      --cri-socket /run/containerd/containerd.sock

You will also need to join the other controllers to this cluster by using the
command provided which includes the ``--control-plane`` flag.  You'll also need
to make sure you add the ``--apiserver-bind-port 16443`` flag or otherwise it
will refuse to join (due to port 6443 being used by the load balancer).

Once that is done, you can proceed to the joining the remainder of the nodes
using the ``kubeadm join`` command that was provided when initializing the
cluster.

After you've completed the installation of the Kubernetes on all of the node,
you can verify that all nodes are present.  It's normal for nodes to be in the
``NotReady`` status due to the fact that the CNI is not present yet::

   $ kubectl get nodes
   NAME   STATUS     ROLES    AGE     VERSION
   ctl1   NotReady   master   17m     v1.18.6
   ctl2   NotReady   master   6m27s   v1.18.6
   ctl3   NotReady   master   5m29s   v1.18.6
   kvm1   NotReady   <none>   18s     v1.18.6
   kvm2   NotReady   <none>   10s     v1.18.6
   kvm3   NotReady   <none>   2s      v1.18.6
