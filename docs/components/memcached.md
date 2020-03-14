# Memcached

Memcached is an in-memory key-value store for small chunks of arbitrary
data (strings, objects) from results of database calls, API calls, or page
rendering.

It's used in OpenStack for a lot of token caching and in other services such
as Nova to minimize load against the database cluster.  This operator allows
you to deploy it and it automatically exposes a single IP address which will
point towards any of the two Mcrouter instances which are pushing data out to
the Memcached instances.

It will also automatically take the total number of megabytes and split it 
across two shards (so setting `megabytes` to `128`) will result in two instances
each with 64 megabytes which are load balanced via Mcrouter.

## Example

```yaml
apiVersion: infrastructure.vexxhost.cloud/v1alpha1
kind: Memcached
metadata:
  name: sample
spec:
  megabytes: 128
```