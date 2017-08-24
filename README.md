Piper
=====

This is a tiny project that sets up a encryption/ decryption proxy between you and IPFS. To the world your files are
random bytes, but to you, your data is private AF. Finally a way to store your data privately on the new internet.

BTW, this is an alpha project. Don't use this with real data. Expect many breaking changes.

### Config

```json
{
  "ApiKey": "super-secret-key-12345",
  "Bind": ":9090",
  "OFBSymmetricKey": "5700826c2d30468d8f6d3361abf9b591"
}
```

### Building

For the CLI tool:
```
$ go get github.com/film42/piper/cmd/piper
```

For the server:
```
$ go get github.com/film42/piper/cmd/piperd
```

### CLI Usage
```
./piper --help
Usage of ./piper:
  -config string
        Path to config file.
  -set string
        Path to file to SET.
```

### Server Usage
```
Usage of ./piperd:
  -config string
        Path to config file.
```

```
$ curl -X POST -H "X-Api-Key: 12345" localhost:9090/set -d "plz don't tell anyone about my gif collection"
QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze%

$ ipfs cat QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze
��5�����[b���K%

$ curl localhost:9090/get/QmZMbfEhWhhoz7ijr33udrmRLB53yzdh2qpyEy9mJ9vUze
plz don't tell anyone about my gif collection%
```

### License

MIT License
