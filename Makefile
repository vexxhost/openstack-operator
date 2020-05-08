images:
	docker build images/horizon -t vexxhost/horizon:latest
	docker build images/keystone -t vexxhost/keystone:latest
	docker build images/ceilometer --target ceilometer-agent-notification -t vexxhost/ceilometer-agent-notification:latest
	docker build images/heat --target heat-api -t vexxhost/heat-api:latest
	docker build images/heat --target heat-api-cfn -t vexxhost/heat-api-cfn:latest
	docker build images/heat --target heat-engine -t vexxhost/heat-engine:latest
	docker build images/mcrouter -t vexxhost/mcrouter:latest
	docker build images/mcrouter-exporter -t vexxhost/mcrouter-exporter:latest
	docker build images/memcached -t vexxhost/memcached:latest
	docker build images/memcached-exporter -t vexxhost/memcached-exporter:latest
	docker build images/rabbitmq -t vexxhost/rabbitmq:latest

