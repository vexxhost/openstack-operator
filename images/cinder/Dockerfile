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
FROM vexxhost/python-base AS cinder-base
RUN mkdir -p /var/lib/cinder && \
    chmod 777 -R /var/lib/cinder

FROM cinder-base AS cinder-api
COPY cinder-wsgi /usr/local/bin/cinder-wsgi
EXPOSE 8776
ENV UWSGI_HTTP_SOCKET=:8776 UWSGI_WSGI_FILE=/usr/local/bin/cinder-wsgi
CMD ["/usr/local/bin/uwsgi", "--ini", "/etc/uwsgi/uwsgi.ini"]

FROM cinder-base AS cinder-volume
COPY cinder-volume /usr/local/bin/cinder-volume
CMD ["/usr/local/bin/cinder-volume"]

FROM cinder-base AS cinder-scheduler
COPY cinder-scheduler /usr/local/bin/cinder-scheduler
CMD ["/usr/local/bin/cinder-scheduler"]