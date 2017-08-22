Personal Space
==============

This is a tiny project that sets up a encryption/ decryption proxy between you and IPFS. To the world your files are
random bytes, but to you, your data is private AF. Sometimes you just need a little personal space, you know? The
project name is a work-in-progress.

### Config

```json
{
  "Bind": ":9090",
  "OFBSymmetricKey": "5700826c2d30468d8f6d3361abf9b591"
}
```

### Building

```
$ go get github.com/film42/personal-space
```

### Usage

```
$ curl -X POST localhost:9090/upload -d "plz don't tell anyone about my gif collection"
QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze%

$ ipfs cat QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze
��5�����[b���K%

$ curl localhost:9090/s/QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze
plz don't tell anyone about my gif collection%
```

### License

MIT License
