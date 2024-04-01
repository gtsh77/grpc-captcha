### GRPC-CAPTCHA [v0.1] [development is currently in progress]

#### requirements
go v1.20+, docker compose v2.x

#### ENV vars
[ENV FILE](.deploy/env/local.env)

#### gRPC captcha schema
[PROTO](pkg/proto/grpc-captcha/grpc-captcha.proto)

#### start storage (non-secure)
```
docker compose -f docker-compose.local.yml up redis -d
```

#### start storage (secure, opt)
```
docker compose -f docker-compose.tls.yml up redis -d
```

#### stop storage, rm volume
```
docker compose -f docker-compose.local.yml down -v
```

#### compile and run

```
make && make run
```

#### compile and run via docker (opt)

```
docker compose -f docker-compose.local.yml up 
```

#### install go-linter (linux, opt)
```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.1
```

#### run linter task
```
make lint
```

#### run autotests
```
make test
```

#### get full env vars list
```
make env
```

#### health check router
```
curl 0.0.0.0:1111/health/check
```

#### operable check router
```
curl 0.0.0.0:1111/health/operable
```

#### prometheus metrics (runtime, http)
```
curl 0.0.0.0:1111/metrics
```

#### redis client local
```
docker run -it --rm --network host redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv
```

#### redis client tls
```
docker run -it --rm --network host -v $(pwd)/.deploy/crt/ca.crt:/ca.crt -v $(pwd)/.deploy/crt/secure.nd.crt:/server.crt -v $(pwd)/.deploy/crt/secure.nd.key:/server.key redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv --tls --cacert /ca.crt --cert /server.crt --key /server.key
```

#### grpc: captcha
```
todo
```