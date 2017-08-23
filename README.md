Personal Space
==============

This is a tiny project that sets up a encryption/ decryption proxy between you and IPFS. To the world your files are
random bytes, but to you, your data is private AF. Sometimes you just need a little personal space, you know? The
project name is a work-in-progress.

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

```
$ go get github.com/film42/personal-space
```

### CLI Usage
```
./personal-space --help
Usage of ./personal-space:
  -config string
        Path to config file.
  -set string
        Path to file to SET.
  -start-server
        Start a gateway server accepting POST / GET requests.
```

You can use `personal-space` as CLI tool for uploading a file. This is IPFS after all, why make the server do both? ;)
The server is not started by default, so you'll need to do that with the `--start-server` option.

### Server Usage

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
