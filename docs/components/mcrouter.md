# Mcrouter

Mcrouter is a memcached protocol router for scaling Memcached deployments. It
is a core component of cache infrastructure at Facebook and Instagram where
Mcrouter handles almost 5 billion requests per second at peak.

OpenStack usually comes down to a halt if one of the Memcached instances don't
respond anymore.  Mcrouter is used to enable high availability and redundancy
so that any Memcached outages will not affect the OpenStack services.

The only two possible reasons that we can have a full system slowdown at the
moment remains:

- All backends (Memcached instances) are all down
- All Mcrouter replicas are down

The first probably means there's a bigger issue in play, the latter will
likely automatically recover by Kubernetes ensuring that replicas come back
up.  Also, due to the fact that the service is exposed as a ClusterIP, it only
takes a single replica to be up for everything to come back to start working
again.

## Example

```yaml
apiVersion: infrastructure.vexxhost.cloud/v1alpha1
kind: Mcrouter
metadata:
  name: sample
spec:
  route: PoolRoute|default
  pools:
    default:
      servers: ['10.0.0.1:11211', '10.0.0.2:11211']
```