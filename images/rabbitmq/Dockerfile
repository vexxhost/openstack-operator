FROM rabbitmq:latest
RUN rabbitmq-plugins enable --offline rabbitmq_prometheus && \
    rabbitmq-plugins enable --offline rabbitmq_management && \
    rabbitmq-plugins enable --offline rabbitmq_peer_discovery_k8s
