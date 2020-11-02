# logstash-exporter

## Build Instructions

Perform the following steps to build a local copy of the logstash-exporter:

1. Clone the repository into your `$GOPATH` or use `go get` to download the repository.

```shell
go get -u github.com/silveiralexandre/logstash-exporter
```

2. From the base directory execute the following:

```shell
$ CGO_ENABLED=0 go build -trimpath .
$
```

3. Execute the compiled binary as shown on the [usage instructions](../docs/usage.md).
