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
RUN apt-get update && \
    apt-get -y install git && \
    git clone https://github.com/openstack/ceilometer /tmp/src
COPY bindep.txt /tmp/src/bindep.txt
RUN assemble

FROM docker.io/opendevorg/python-base AS ceilometer
COPY --from=builder /output/ /output
RUN /output/install-from-bindep

FROM ceilometer AS ceilometer-agent-notification
CMD ["/usr/local/bin/ceilometer-agent-notification"]
