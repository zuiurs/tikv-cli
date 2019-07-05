# tikv-cli

Simple command-line client for TiKV.

## Usage

```
$ tikv-cli -h 0.0.0.0:23791
INFO[0000] [pd] create pd client with endpoints [0.0.0.0:23791]
INFO[0000] [pd] leader switches to: http://172.17.0.1:23791, previous:
INFO[0000] [pd] init cluster id 6709365086357114166
> get a

> put a b
successed!
> get a
b
> delete a
successed!
> get a

```

## License

MIT
