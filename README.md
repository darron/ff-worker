# ff-worker

Cloudflare Worker way to run [ff](https://github.com/darron/ff)

Based on template:

## worker-template-tinygo

- A template for starting a Cloudflare Worker project with tinygo.
- This template uses [`workers`](https://github.com/syumai/workers) package to run an HTTP server.

## Usage

- `main.go` includes simple HTTP server implementation. Feel free to edit this code and implement your own HTTP server.

## Requirements

- Node.js
- [wrangler](https://developers.cloudflare.com/workers/wrangler/)
  - just run `npm install -g wrangler`
- tinygo 0.29.0 or later

## Getting Started

* If not already installed, please install the [gonew](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) command.

```console
go install golang.org/x/tools/cmd/gonew@latest
```

* Create a new project using this template.
  - Second argument passed to `gonew` is a module path of your new app.

```console
gonew github.com/syumai/workers/_templates/cloudflare/worker-tinygo your.module/my-app # e.g. github.com/syumai/my-app
cd my-app
go mod tidy
make dev # start running dev server
curl http://localhost:8787/hello # outputs "Hello!"
```

- To change worker name, please edit `name` property in `wrangler.toml`.

## Development

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy # deploy worker
```

### Testing dev server

- Just send HTTP request using some tools like curl.

```
$ curl http://localhost:8787/hello
Hello!
```

```
$ curl -X POST -d "test message" http://localhost:8787/echo
test message
```
