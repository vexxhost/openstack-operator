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
FROM vexxhost/python-base AS nova-base
RUN mkdir -p /var/lib/nova && \
    chmod 777 -R /var/lib/nova

FROM nova-base AS nova-compute-api
COPY nova-api-wsgi /usr/local/bin/nova-api-wsgi
EXPOSE 8774
ENV UWSGI_HTTP_SOCKET=:8774 UWSGI_WSGI_FILE=/usr/local/bin/nova-api-wsgi
CMD ["/usr/local/bin/uwsgi", "--ini", "/etc/uwsgi/uwsgi.ini"]

FROM nova-base AS nova-metadata-api
COPY nova-metadata-wsgi /usr/local/bin/nova-metadata-wsgi
EXPOSE 8775
ENV UWSGI_HTTP_SOCKET=:8775 UWSGI_WSGI_FILE=/usr/local/bin/nova-metadata-wsgi
CMD ["/usr/local/bin/uwsgi", "--ini", "/etc/uwsgi/uwsgi.ini"]

FROM nova-base AS nova-conductor
COPY nova-init /usr/local/bin/nova-init
COPY nova-conductor /usr/local/bin/nova-conductor
CMD ["/usr/local/bin/nova-conductor"]

FROM nova-base AS nova-scheduler
COPY nova-scheduler /usr/local/bin/nova-scheduler
CMD ["/usr/local/bin/nova-scheduler"]

FROM debian:10-slim AS novnc-builder
ADD https://github.com/novnc/noVNC/archive/v1.2.0.tar.gz novnc.tgz
RUN tar -xvzf novnc.tgz

FROM nova-base AS nova-novncproxy
COPY nova-novncproxy /usr/local/bin/nova-novncproxy
COPY --from=novnc-builder noVNC-1.2.0/ /usr/share/novnc
EXPOSE 6080
CMD ["/usr/local/bin/nova-novncproxy"]

FROM nova-base AS nova-compute
COPY nova-compute /usr/local/bin/nova-compute
CMD ["/usr/local/bin/nova-compute"]