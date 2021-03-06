#!/bin/bash
#
# Copyright 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

function copy_minikube_config {
	mkdir ~stack/.kube

	sudo cp ~zuul/.kube/config ~stack/.kube/config
	sudo cp ~zuul/.minikube/ca.crt ~stack/.kube/ca.crt
	sudo cp ~zuul/.minikube/profiles/minikube/client.crt ~stack/.kube/client.crt
	sudo cp ~zuul/.minikube/profiles/minikube/client.key ~stack/.kube/client.key
	sudo chown -Rv stack:stack ~stack/.kube

	sed -i s%/home/zuul/.minikube/profiles/minikube%/opt/stack/.kube% ~/.kube/config
	sed -i s%/home/zuul/.minikube/ca.crt%/opt/stack/.kube/ca.crt% ~/.kube/config

	kubectl cluster-info
}

if [[ "$1" == "stack" && "$2" == "pre-install" ]]; then
	copy_minikube_config

elif [[ "$1" == "stack" && "$2" == "install" ]]; then
	:

elif [[ "$1" == "stack" && "$2" == "post-config" ]]; then
	# After ceph devstack plugin
    kubectl create secret generic glance-config -n openstack \
    --from-file=/etc/glance/glance-api.conf \
    --from-file=/etc/glance/glance-api-paste.ini

    kubectl create secret generic cinder-config -n openstack \
    --from-file=/etc/cinder/cinder.conf \
    --from-file=/etc/cinder/api-paste.ini \
    --from-file=/etc/cinder/rootwrap.conf \
    --from-file=/etc/cinder/resource_filters.json

	# NOTE(Alex): Permissions here are bad but it's temporary so we don't care as much.
	sudo chmod -Rv 777 /etc/ceph
    kubectl create secret generic ceph-config -n openstack \
    --from-file=/etc/ceph/ceph.conf \
	--from-file=/etc/ceph/ceph.client.cinder.keyring \
	--from-file=/etc/ceph/ceph.client.glance.keyring

	# NOTE(Alex): Create nova compute conf to include placement and libvirt config
	create_nova_compute_conf
	# NOTE(Alex) To include create_nova_conf_neutron and barbican hack config
	kubectl create secret generic nova-config -n openstack \
	--from-file=/etc/nova/nova.conf \
	--from-file=/etc/nova/nova-cpu.conf \
	--from-file=/etc/nova/nova_cell1.conf \
	--from-file=/etc/nova/api-paste.ini

elif [[ "$1" == "stack" && "$2" == "extra" ]]; then
	:

elif [[ "$1" == "stack" && "$2" == "test-config" ]]; then
	# Horizon dashboard Url in tempest_horizon
	if is_service_enabled tempest; then
		local ip=$(get_kubernetes_service_ip horizon)
		iniset $TEMPEST_CONFIG dashboard dashboard_url http://$ip
	fi
fi

if [[ "$1" == "unstack" ]]; then
	:
fi

if [[ "$1" == "clean" ]]; then
	:
fi
