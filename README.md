# Selenoid API Proxy
Tiny web server, which runs in front of Selenoid API and provides token-based authorization.

## Usage
* Delete video
```shell
curl -X DELETE -H "X-Token:token-here" https://grid-api.kyberorg.io/video/yalsee-dev-130821-1012
```

## Run
* Minimal configuration example
```shell
selenoid-api-proxy --token myCustomTokenHere
```

* Fully customized params example
```shell
selenoid-api-proxy --port 4443 --selenoid.api.url=http://localhost:4444  --token myCustomTokenHere
```

or same with env vars
```shell
SELENOID_API_URL=http://localhost:4444 PORT=4443 TOKEN=myCustomTokenHere selenoid-api-proxy
```

### Help
```shell
selenoid-api-proxy --help
```

### Systemd Daemon
```unit file (systemd)
[Unit]
Description=Selenoid API Proxy
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=5s
Environment="TOKEN=my-token-here"
WorkingDirectory=/srv/selenoid-api-proxy
ExecStart=/srv/selenoid-api-proxy/selenoid-api-proxy \
    --port=4443 \
    --selenoid.api.url=http://localhost:4444

SyslogIdentifier=selenoid-api-proxy

[Install]
WantedBy=multi-user.target
```

## Build
```shell
make binary
```

or without `make`
```shell
 CGO_ENABLED=0 go build github.com/kyberorg/selenoid-api-proxy/cmd/selenoid-api-proxy
```

