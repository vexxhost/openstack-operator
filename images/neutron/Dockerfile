# Copyright (c) 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM vexxhost/python-builder as builder
FROM vexxhost/python-base AS neutron-base
RUN mkdir -p /var/lib/neutron && \
    chmod 777 -R /var/lib/neutron

FROM neutron-base AS neutron-api
COPY neutron-api /usr/local/bin/neutron-api
EXPOSE 9696
ENV UWSGI_HTTP_SOCKET=:9696 UWSGI_WSGI_FILE=/usr/local/bin/neutron-api
CMD ["/usr/local/bin/uwsgi", "--ini", "/etc/uwsgi/uwsgi.ini"]

FROM neutron-base AS neutron-rpc-server
COPY neutron-rpc-server /usr/local/bin/neutron-rpc-server
CMD ["/usr/local/bin/neutron-rpc-server"]

FROM neutron-base AS neutron-openvswitch-agent
COPY neutron-openvswitch-agent /usr/local/bin/neutron-openvswitch-agent
CMD ["/usr/local/bin/neutron-openvswitch-agent", "--config-file", "/etc/neutron/neutron.conf", "--config-file", "/etc/neutron/plugins/ml2/ml2_conf.ini"]

FROM neutron-base AS neutron-l3-agent
COPY neutron-l3-agent /usr/local/bin/neutron-l3-agent
CMD ["/usr/local/bin/neutron-l3-agent", "--config-file", "/etc/neutron/neutron.conf", "--config-file", "/etc/neutron/l3_agent.ini", "--config-file", "/etc/neutron/plugins/ml2/ml2_conf.ini"]

FROM neutron-base AS neutron-dhcp-agent
COPY neutron-dhcp-agent /usr/local/bin/neutron-dhcp-agent
CMD ["/usr/local/bin/neutron-dhcp-agent", "--config-file", "/etc/neutron/neutron.conf", "--config-file", "/etc/neutron/dhcp_agent.ini", "--config-file", "/etc/neutron/plugins/ml2/ml2_conf.ini"]

FROM neutron-base AS neutron-metadata-agent
COPY neutron-metadata-agent /usr/local/bin/neutron-metadata-agent
CMD ["/usr/local/bin/neutron-metadata-agent", "--config-file", "/etc/neutron/neutron.conf", "--config-file", "/etc/neutron/metadata_agent.ini", "--config-file", "/etc/neutron/plugins/ml2/ml2_conf.ini"]