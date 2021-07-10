# README

## Env Parameters:
```
export corpid=
export corpsecret=
export token=
export debug=
export oapi_server=oapi.dingtalk.com
export timezone="Asia/Shanghai"
```

## Build binary for R2S router
```
DOCKER_BUILDKIT=1 docker build --file Dockerfile --output out .
```