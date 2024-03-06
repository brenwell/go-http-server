# Go-HTTP-Server

This is my test server for Kubernetes.

It uses `buildx` to build multi architecture. This is how you set it up locally.

```sh
function setup_builder {
  docker buildx create --name multiarch-builder --driver docker-container --use
  docker buildx inspect --bootstrap
}
```

And this is how you build it.

```sh
function build {
  docker buildx build --platform linux/amd64,linux/arm64 -t brenwell/go-http-server:latest -f app/Dockerfile --push app/
}
```
