#### redis client local
```
docker run -it --rm --network host redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv
```

#### redis client tls
```
docker run -it --rm --network host -v $(pwd)/.deploy/crt/ca.crt:/ca.crt -v $(pwd)/.deploy/crt/secure.nd.crt:/server.crt -v $(pwd)/.deploy/crt/secure.nd.key:/server.key redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv --tls --cacert /ca.crt --cert /server.crt --key /server.key
```
