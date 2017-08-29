# Autotest library

A golang library to create automatic tests/integration tests easily.

## Get started

```
go get github.com/osechet/autotest
```

### Run unit tests

```
go test -timeout 30s github.com/osechet/autotest/...
```

### Run integration tests

```
go test -timeout 30s -tags=integration github.com/osechet/autotest
```

### Lint

```
go tool vet -v src/github.com/osechet/autotest
```

#### Coverage

```
go test -coverprofile=coverage.out github.com/osechet/autotest
go tool cover -html=coverage.out
```
