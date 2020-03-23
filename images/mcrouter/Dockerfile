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
