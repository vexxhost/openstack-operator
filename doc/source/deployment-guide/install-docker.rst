Install Docker
--------------
Docker is used by many different components of the underlying infrastructure,
so it must be installed before anything to bootstrap the system.

It must be installed on all the machines that you intend to manage using the
OpenStack operator.  It will also be used to deploy infrastructure components
such as the virtual IP and Ceph.

.. tabs::

   .. code-tab:: console Debian

      $ apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
      $ curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
      $ sudo add-apt-repository "deb https://download.docker.com/linux/debian $(lsb_release -cs) stable"
      $ apt-get update
      $ apt-get install -y docker-ce
      $ apt-mark hold docker-ce

