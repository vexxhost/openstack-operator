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
RUN apt update && apt -y install git
ADD https://releases.openstack.org/constraints/upper/ussuri /tmp/src/upper-constraints.txt
RUN assemble git+https://opendev.org/openstack/glance@stable/ussuri \
             boto3 \
             msgpack===0.6.1 \
             PyMySQL \
	         python-memcached \
             python-swiftclient

FROM docker.io/opendevorg/uwsgi-base as glance-api
COPY --from=builder /output/ /output
RUN apt update && apt -y install git && \
    /output/install-from-bindep && \
    apt -y purge git && apt clean all
EXPOSE 9292
ENV UWSGI_HTTP_SOCKET=:9292 UWSGI_WSGI_FILE=/usr/local/bin/glance-wsgi-api
CMD ["/usr/local/bin/uwsgi","--ini","/etc/uwsgi/uwsgi.ini"]