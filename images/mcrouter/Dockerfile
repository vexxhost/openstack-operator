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

FROM ubuntu:bionic
RUN groupadd -r mcrouter \
  && useradd -r -g mcrouter mcrouter
RUN apt update && \
    apt install -y --no-install-recommends ca-certificates wget gnupg && \
    wget -O - https://facebook.github.io/mcrouter/debrepo/bionic/PUBLIC.KEY | apt-key add && \
    echo "deb https://facebook.github.io/mcrouter/debrepo/bionic bionic contrib" >> /etc/apt/sources.list && \
    apt update && \
    apt install -y mcrouter && \
    apt remove -y wget gnupg && \
    apt autoremove -y && \
    apt clean all
RUN chown -R mcrouter:mcrouter /var/spool/mcrouter
RUN chown -R mcrouter:mcrouter /var/mcrouter
RUN chown -R mcrouter:mcrouter /usr/bin/mcrouter
USER mcrouter
ENTRYPOINT ["/usr/bin/mcrouter"]
