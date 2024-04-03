### GRPC-CAPTCHA [v0.2] [development is currently in progress]

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

#### grpc: generate new captcha image (insecure)
```
grpcurl -H 'x-api-key: 1ace3bed-3aaf-4642-adb1-d63aef85895f' -plaintext -import-path pkg/proto/grpc-captcha -proto grpc-captcha.proto 0.0.0.0:2222 werkstatt.captcha.CaptchaService.Generate |jq

{
  "id": "acf26399-0aa3-4fea-89ef-495476315998",
  "data": "iVBORw0KGgoAAAANSUhEUgAAALQAAABQCAMAAACHxq+UAAAAP1BMVEUAAAA5QBDe5bWzuord5LSwt4fZ4LBzekq5wJDHzp7g57dzekp6gVEyOQlUWyuMk2NSWSmZoHDHzp6Ij1+xuIidZE4xAAAAAXRSTlMAQObYZgAAA29JREFUeJzsmt2yqyoMx5Npay9qp07f/2HPVBGSEL4R19ljerF3K4Qff0KIuuCyyy677LLLLrvsVMOzASoMERHwf4aOK/bfJ5+myf7f4h4Bbr293+9GV9MkqS3w+m8f21cQNuZWagWaKh2lfuThEikOgQbjOq60+fXxSFAzYB36+XxWUnvMbkwNxbZJMXNg64wz11HzgRIRARlRY1uyHaK26QKd3HsF+zMneXZSOk4U2J9oP17ThPVgTkpNlH6yTnSrZXnqafGxiNJWIwRtu41kjq/qekUojdJc26FKB7aPQWJXUPnAeKFZ7aExUx4hrvs2FBp47eFd0g4g5ZubyjIIO1DlIQGi85B5jk1sWZZxke2fIk56ssXkahBmE/3L0DjhSCTIvV95h31eoCXvUdAGlOhLlZablYS0euRk2dRGDUo20PRHQhrJ3ZnMUxs1OQJZjhBJjtAxfHvDUkjdxsyovmQe7CIQefdp2UNGS55Hm5Pr+/3an9TINfxKdh/MbAdFBK40okpt2opyb/jjCG8v2d3nItcGhGnKDs0TmH2lnJw8SvaL7tJ50DIDsLPSr+9MFy1ZVg3eAD3PM81hIn/4IaSf+lVjV/ecf8Y2HpA6lijOp4rYqnR+X3nbaqB5hgCub0jrltxR0lnZeEJpYPThEbAlSxdOmEchmpgGVen4GAOZ+Q0AWWT+nE/UFNp9pVuYGuii1MHuXOTuY0+xfUKFuUbu4iILQL4SkAOjWmWo5XUlN9ZlaOTjsVPObUa/i/jOFqto7Apmr4hTr/k9ZKnEe3Rhjnvx48J6DcY0qaA04XOwk+2SbnQPwewhg6qMJi/TpbJRwEUkT8d2XjK0Mw+jeBYNDrLJeLvd/Ed5YeYEtVmbrBAKNYy+llsX/vYzKS0J6ogMxOm7ghlsYKN/nx33sUL7FXVyk7AnxutbIwSXQ8vyIhXLHYkRH1Lp3NHIkYuIbzdUYW4mSttFTvv4Ca2+b8ka0hW1ez6oKlDc9ItnXWF6hVvrC0tuGNqtz1jYcD952b9q97MBKux+P5Y6+YcWVXY48yHUh9og6Lmrt2ro9djNtXnuTl3TbavKsq0vc60VQrdbj8UaDb09AGy18UJ33hcDTIM+5VVGkanMf57as0OgX70dCtugP119vl4DqOHz6Uzd1VvI2pj/CwAA//9L/w3JNJ6cGwAAAABJRU5ErkJggg=="
}
```

#### grpc: validate captcha (insecure)
```
grpcurl -H 'x-api-key: 1ace3bed-3aaf-4642-adb1-d63aef85895f' -plaintext -import-path pkg/proto/grpc-captcha -proto grpc-captcha.proto -d '{"id":"acf26399-0aa3-4fea-89ef-495476315998", "code": "5485"}' 0.0.0.0:2222 werkstatt.captcha.CaptchaService.Verify |jq

{}
```

#### render png
paste data field value into 
https://onlinepngtools.com/convert-base64-to-png