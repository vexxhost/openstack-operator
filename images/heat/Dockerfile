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
FROM vexxhost/python-builder AS builder
FROM vexxhost/python-base AS heat-api-base

FROM heat-api-base AS heat-api
COPY heat-wsgi-api /usr/local/bin/heat-wsgi-api
EXPOSE 8004
ENV UWSGI_HTTP_SOCKET=:8004 UWSGI_WSGI_FILE=/usr/local/bin/heat-wsgi-api
CMD ["/usr/local/bin/uwsgi","--ini","/etc/uwsgi/uwsgi.ini"]

FROM heat-api-base AS heat-api-cfn
COPY heat-wsgi-api-cfn /usr/local/bin/heat-wsgi-api-cfn
EXPOSE 8000
ENV UWSGI_HTTP_SOCKET=:8000 UWSGI_WSGI_FILE=/usr/local/bin/heat-wsgi-api-cfn
CMD ["/usr/local/bin/uwsgi","--ini","/etc/uwsgi/uwsgi.ini"]

FROM heat-api-base AS heat-engine
COPY heat-engine /usr/local/bin/heat-engine
RUN ln -s /usr/local/etc/heat /etc/heat
CMD ["/usr/local/bin/heat-engine"]
