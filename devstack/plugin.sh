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
	:

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
