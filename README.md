### GRPC-CAPTCHA [v1.0.0] (12 Apr 2024)
### INFO

#### requirements
go v1.20+, redis v7.x, docker compose v2.x

#### description
captcha code sequence generator, image render and verify code service with gRPC interfaces<br>
based on lite (render, random) version of https://github.com/dchest/captcha

##### example
[![image](https://i.postimg.cc/rshMqnnR/captcha4d.png)](https://i.postimg.cc/rshMqnnR/captcha4d.png)
>180x80px 4 digits png image

#### ENV vars
[ENV FILE](.deploy/env/local.env)

#### gRPC service schema
[PROTO](pkg/proto/grpc-captcha/grpc-captcha.proto)

#### Multi platform binaries
[LIST](bin)

### SIMPLE USAGE (NON-TLS)

#### start via docker image (specify own redis 7.x instance)
```
docker run --rm -it -p 1111:1111 -p 2222:2222 --env-file=.deploy/env/local.env -e CAPTCHA_REDIS_HOST=0.0.0.0 -e CAPTCHA_REDIS_PORT=6379 -e CAPTCHA_REDIS_DB=0 -e CAPTCHA_REDIS_PASS=YQ3dvPx3fVzv gtsh77workshop/grpc-captcha:v1.0.0
```

#### OR start via docker compose with redis
```
docker compose -f docker-compose.local.yml up
```

#### OR export env settings and use binaries depends on your platform
```
export $(grep -v '^#' ./deploy/env/local.env | xargs)
./bin/grpc-captcha_linux_amd64
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
grpcurl -H 'x-api-key: 1ace3bed-3aaf-4642-adb1-d63aef85895f' -plaintext -import-path pkg/proto/grpc-captcha -proto grpc-captcha.proto -d '{"id":"acf26399-0aa3-4fea-89ef-495476315998", "otp": "5485"}' 0.0.0.0:2222 werkstatt.captcha.CaptchaService.Verify |jq
{}
```

#### render png
paste data field value into 
https://onlinepngtools.com/convert-base64-to-png

#### redis client (opt)
```
docker run -it --rm --network host redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv
```

### SECURED USAGE (TLSv1.2/v1.3)

For most cases we recommend simple use local proxy (eg `nginx`) with TLS enabled which upstreams non-tls captcha service or even simpler you can use `k8s` which already have its own TLS proxy/balancer

If in some cases there is no possibility to have a local proxy **we prepared this full operable example for direct interaction with service's interfaces through TLS**

1. load our test self-signed crt/keys data into env variables or use own (crt, key and ca crt) [using data instead of file path helps with k8s pod implementation, so you need to load them in advance]
    ```
    export CAPTCHA_HTTP_TLS_CRT_DATA=`cat .deploy/crt/secure.nd.crt` CAPTCHA_GRPC_TLS_CRT_DATA=`cat .deploy/crt/secure.nd.crt` CAPTCHA_REDIS_TLS_CRT_DATA=`cat .deploy/crt/secure.nd.crt` CAPTCHA_HTTP_TLS_KEY_DATA=`cat .deploy/crt/secure.nd.key` CAPTCHA_GRPC_TLS_KEY_DATA=`cat .deploy/crt/secure.nd.key` CAPTCHA_REDIS_TLS_KEY_DATA=`cat .deploy/crt/secure.nd.key` CAPTCHA_HTTP_TLS_CRT_CA_DATA=`cat .deploy/crt/ca.crt`  CAPTCHA_GRPC_TLS_CRT_CA_DATA=`cat .deploy/crt/ca.crt` CAPTCHA_REDIS_TLS_CRT_CA_DATA=`cat .deploy/crt/ca.crt`
    ```

2. start via docker compose (tls variant)
    ```
    docker compose -f docker-compose.tls.yml up
    ```

3. add crt to trusted (only if self-signed) to be able to interact with grpc interface via tls
    ```
      apt install ca-certificates
      cp .deploy/crt/secure.nd.crt /usr/local/share/ca-certificates/.
      update-ca-certificates
    ```

4. set local domain which corresponds to your crt's CN (opt) [below example for test self-signed crt/key provided within repo]
    ```
    vi /etc/hosts
    captcha.secure.nd 127.0.0.1
    ```

#### grpc: generate new captcha image (secure)
```
grpcurl -H 'x-api-key: 1ace3bed-3aaf-4642-adb1-d63aef85895f' -import-path pkg/proto/grpc-captcha -proto grpc-captcha.proto captcha.secure.nd:2222 werkstatt.captcha.CaptchaService.Generate |jq
{
  "id": "acf26399-0aa3-4fea-89ef-495476315998",
  "data": "iVBORw0KGgoAAAANSUhEUgAAALQAAABQCAMAAACHxq+UAAAAP1BMVEUAAAA5QBDe5bWzuord5LSwt4fZ4LBzekq5wJDHzp7g57dzekp6gVEyOQlUWyuMk2NSWSmZoHDHzp6Ij1+xuIidZE4xAAAAAXRSTlMAQObYZgAAA29JREFUeJzsmt2yqyoMx5Npay9qp07f/2HPVBGSEL4R19ljerF3K4Qff0KIuuCyyy677LLLLrvsVMOzASoMERHwf4aOK/bfJ5+myf7f4h4Bbr293+9GV9MkqS3w+m8f21cQNuZWagWaKh2lfuThEikOgQbjOq60+fXxSFAzYB36+XxWUnvMbkwNxbZJMXNg64wz11HzgRIRARlRY1uyHaK26QKd3HsF+zMneXZSOk4U2J9oP17ThPVgTkpNlH6yTnSrZXnqafGxiNJWIwRtu41kjq/qekUojdJc26FKB7aPQWJXUPnAeKFZ7aExUx4hrvs2FBp47eFd0g4g5ZubyjIIO1DlIQGi85B5jk1sWZZxke2fIk56ssXkahBmE/3L0DjhSCTIvV95h31eoCXvUdAGlOhLlZablYS0euRk2dRGDUo20PRHQhrJ3ZnMUxs1OQJZjhBJjtAxfHvDUkjdxsyovmQe7CIQefdp2UNGS55Hm5Pr+/3an9TINfxKdh/MbAdFBK40okpt2opyb/jjCG8v2d3nItcGhGnKDs0TmH2lnJw8SvaL7tJ50DIDsLPSr+9MFy1ZVg3eAD3PM81hIn/4IaSf+lVjV/ecf8Y2HpA6lijOp4rYqnR+X3nbaqB5hgCub0jrltxR0lnZeEJpYPThEbAlSxdOmEchmpgGVen4GAOZ+Q0AWWT+nE/UFNp9pVuYGuii1MHuXOTuY0+xfUKFuUbu4iILQL4SkAOjWmWo5XUlN9ZlaOTjsVPObUa/i/jOFqto7Apmr4hTr/k9ZKnEe3Rhjnvx48J6DcY0qaA04XOwk+2SbnQPwewhg6qMJi/TpbJRwEUkT8d2XjK0Mw+jeBYNDrLJeLvd/Ed5YeYEtVmbrBAKNYy+llsX/vYzKS0J6ogMxOm7ghlsYKN/nx33sUL7FXVyk7AnxutbIwSXQ8vyIhXLHYkRH1Lp3NHIkYuIbzdUYW4mSttFTvv4Ca2+b8ka0hW1ez6oKlDc9ItnXWF6hVvrC0tuGNqtz1jYcD952b9q97MBKux+P5Y6+YcWVXY48yHUh9og6Lmrt2ro9djNtXnuTl3TbavKsq0vc60VQrdbj8UaDb09AGy18UJ33hcDTIM+5VVGkanMf57as0OgX70dCtugP119vl4DqOHz6Uzd1VvI2pj/CwAA//9L/w3JNJ6cGwAAAABJRU5ErkJggg=="
}
```

#### grpc: validate captcha (secure)
```
grpcurl -H 'x-api-key: 1ace3bed-3aaf-4642-adb1-d63aef85895f' -import-path pkg/proto/grpc-captcha -proto grpc-captcha.proto -d '{"id":"acf26399-0aa3-4fea-89ef-495476315998", "code": "5485"}' captcha.secure.nd:2222 werkstatt.captcha.CaptchaService.Verify |jq
{}
```

#### common http routes (secure)
```
use curl -k https://host... instead of curl host...
```

#### redis client with tls support (opt)
```
docker run -it --rm --network host -v $(pwd)/.deploy/crt/ca.crt:/ca.crt -v $(pwd)/.deploy/crt/secure.nd.crt:/server.crt -v $(pwd)/.deploy/crt/secure.nd.key:/server.key redis:7.2-alpine redis-cli -p 6379 -a YQ3dvPx3fVzv --tls --cacert /ca.crt --cert /server.crt --key /server.key
```

#### render png
paste data field value into 
https://onlinepngtools.com/convert-base64-to-png

### DEV INFO

#### design
##### method CaptchaService.Generate(empty):
Generate new sequence code of [0-9] digits in quantity of `CAPTCHA_RENDER_DIG_CNT`, render png image `CAPTCHA_RENDER_WIDTH` x `CAPTCHA_RENDER_HEIGHT` px, store code sequence in redis db with TTL `CAPTCHA_RENDER_TTL` and return base64 encoding of image and uniqique captcha `id` to identify the sequence (`uuidv4`)

##### method CaptchaService.Verify(id, code):
compare code sequences in payload and value stored in redis by provided key (`id`), if no such key ret `NotFound`, if values are not equal ret `FailedPrecondition`, if eq ret `OK` and remove key from redis

#### compile and run on your host

```
make && make run
```

#### install go-linter (linux, opt)
```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.1
```

#### run linter task (opt)
```
make lint
```

#### run autotests (opt)
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


define B = 
go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(APP_OUT) $(APP_MAIN)
endef