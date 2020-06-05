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

# Save trace setting
XTRACE=$(set +o | grep xtrace)
set -o xtrace

echo_summary "magnum's plugin.sh was called..."
source $DEST/magnum/devstack/lib/magnum
(set -o posix; set)

if is_service_enabled magnum-api magnum-cond; then
    if [[ "$1" == "stack" && "$2" == "install" ]]; then
        echo_summary "Installing magnum"
        install_magnum

        MAGNUM_GUEST_IMAGE_URL=${MAGNUM_GUEST_IMAGE_URL:-"https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/31.20200323.3.2/x86_64/fedora-coreos-31.20200323.3.2-openstack.x86_64.qcow2.xz"}
        IMAGE_URLS+=",${MAGNUM_GUEST_IMAGE_URL}"

        LIBS_FROM_GIT="${LIBS_FROM_GIT},python-magnumclient"

        install_magnumclient
        cleanup_magnum
    elif [[ "$1" == "stack" && "$2" == "post-config" ]]; then
        echo_summary "Configuring magnum"
        configure_magnum

        if is_service_enabled key; then
            create_magnum_accounts
        fi

    elif [[ "$1" == "stack" && "$2" == "extra" ]]; then
        # Initialize magnum
        init_magnum
        magnum_register_image
        magnum_configure_flavor

        # Start the magnum API and magnum taskmgr components
        echo_summary "Starting magnum"
        start_magnum

        configure_iptables_magnum
        configure_apache_magnum
    fi

    if [[ "$1" == "unstack" ]]; then
        stop_magnum
    fi

    if [[ "$1" == "clean" ]]; then
        cleanup_magnum
    fi
fi

# Restore xtrace
$XTRACE
