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

FROM vexxhost/python-base AS magnum-api
COPY magnum-api-wsgi /usr/local/bin/magnum-api-wsgi
EXPOSE 9511
ENV UWSGI_HTTP_SOCKET=:9511 UWSGI_WSGI_FILE=/usr/local/bin/magnum-api-wsgi
CMD ["/usr/local/bin/uwsgi","--ini","/etc/uwsgi/uwsgi.ini"]

FROM vexxhost/python-base AS magnum-conductor
COPY magnum-conductor /usr/local/bin/magnum-conductor
CMD ["/usr/local/bin/magnum-conductor"]
