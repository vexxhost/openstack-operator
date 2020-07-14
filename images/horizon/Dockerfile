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
RUN assemble git+https://opendev.org/openstack/horizon@stable/ussuri \
             git+https://opendev.org/openstack/heat-dashboard@stable/ussuri \
             git+https://opendev.org/openstack/designate-dashboard@stable/ussuri \
             git+https://opendev.org/openstack/neutron-vpnaas-dashboard@stable/ussuri \
             git+https://opendev.org/openstack/octavia-dashboard@stable/ussuri \
             git+https://opendev.org/openstack/sahara-dashboard@stable/ussuri \
             git+https://opendev.org/openstack/magnum-ui@stable/ussuri \
             python-memcached

FROM docker.io/opendevorg/uwsgi-base
COPY --from=builder /output/ /output
RUN apt update && apt -y install git && \
    /output/install-from-bindep && \
    apt -y purge git && apt clean all
WORKDIR /usr/local/lib/python3.7/site-packages/openstack_dashboard
RUN ln -s /etc/horizon/local_settings.py local/local_settings.py && \
    cp ../designatedashboard/enabled/*.py enabled/ && \
    cp ../heat_dashboard/enabled/*.py enabled/ && \
    cp ../magnum_ui/enabled/*.py enabled/ && \
    cp ../neutron_vpnaas_dashboard/enabled/*.py enabled/ && \
    cp ../octavia_dashboard/enabled/*.py enabled/ && \
    cp ../sahara_dashboard/enabled/*.py enabled/
COPY manage.py .
RUN python manage.py compilemessages && \
    python manage.py collectstatic --no-input && \
    python manage.py compress --force && \
    chown 1001 -R local/ ../static
COPY *.svg ../static/dashboard/img/
COPY *.ico ../static/dashboard/img/
EXPOSE 8000
ENV UWSGI_HTTP_SOCKET=:8000 \
    UWSGI_WSGI_FILE=/usr/local/lib/python3.7/site-packages/openstack_dashboard/wsgi.py \
    UWSGI_CHECK_STATIC=/usr/local/lib/python3.7/site-packages/static/ \
    UWSGI_STATIC_MAP="/static=/usr/local/lib/python3.7/site-packages/static/" \
    UWSGI_MIME_FILE="/etc/mime.types"
CMD ["/usr/local/bin/uwsgi","--ini","/etc/uwsgi/uwsgi.ini"]
