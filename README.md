# tikv-cli

Simple command-line client for TiKV.

## Usage

- One shot

```
$ tikv-cli -h tikv-pd:2379 put hello world
INFO[0000] [pd] create pd client with endpoints [tikv-pd:2379]
INFO[0000] [pd] leader switches to: http://tikv-pd-2.tikv-pd-peer.default.svc:2379, previous:
INFO[0000] [pd] init cluster id 6729357499514264512
successed!

$ tikv-cli -h tikv-pd:2379 get hello
INFO[0000] [pd] create pd client with endpoints [tikv-pd:2379]
INFO[0000] [pd] leader switches to: http://tikv-pd-2.tikv-pd-peer.default.svc:2379, previous:
INFO[0000] [pd] init cluster id 6729357499514264512
world
```

- Interactive shell

```
$ tikv-cli -h tikv-pd:2379
INFO[0000] [pd] create pd client with endpoints [tikv-pd:2379]
INFO[0000] [pd] leader switches to: http://tikv-pd-2.tikv-pd-peer.default.svc:2379, previous:
INFO[0000] [pd] init cluster id 6729357499514264512
> put hello world
successed!
> get hello
world
>
```

## Installation

Use `go get`.

```
go get -u github.com/zuiurs/tikv-cli
```

## Docker Image

- https://hub.docker.com/r/zuiurs/tikv-cli

## Run on Kubernetes

```
kubectl run tikv-cli \
  --image zuiurs/tikv-cli \
  --image-pull-policy Always \
  --restart Never \
  -it --rm -- \
  /bin/sh
```

## License

MIT
