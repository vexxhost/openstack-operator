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

FROM docker.io/opendevorg/python-builder as builder
COPY bindep.txt /tmp/src/bindep.txt
ADD https://releases.openstack.org/constraints/upper/ussuri /tmp/src/upper-constraints.txt
RUN assemble openstack-heat==14.0.0 \
    python-memcached \
    PyMySQL

FROM docker.io/opendevorg/uwsgi-base AS heat-api-base
COPY --from=builder /output/ /output
RUN /output/install-from-bindep

FROM heat-api-base AS heat-api
EXPOSE 8004
ENV UWSGI_HTTP_SOCKET=:8004 UWSGI_WSGI_FILE=/usr/local/bin/heat-wsgi-api

FROM heat-api-base AS heat-api-cfn
EXPOSE 8000
ENV UWSGI_HTTP_SOCKET=:8000 UWSGI_WSGI_FILE=/usr/local/bin/heat-wsgi-api-cfn

FROM docker.io/opendevorg/python-base AS heat-engine
COPY --from=builder /output/ /output
RUN /output/install-from-bindep
RUN ln -s /usr/local/etc/heat /etc/heat
CMD ["/usr/local/bin/heat-engine"]
