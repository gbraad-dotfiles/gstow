# Tapewrap

### info

Build actions for [`tapewrap`](https://github.com/ducttape-infra/tapewrap), a cross-platform
GNU stow replacement written in Go.

### config
```ini
[devenv]
    name="gobuild"
    from="gofedora"
```

### vars
```sh
# hardcode /home/gbraad instead of ~/
COMPILE_REPO_LOCAL=$(eval echo "${PWD}")
```

### local-build
```sh evaluate
make
```

### make
One-shot ephemeral build for local arch.

```sh evaluate
devenv ${DEVENV_FROM} ephemeral -c "cd ${COMPILE_REPO_LOCAL} && make"
```

### test
```sh
devenv ${DEVENV_FROM} ephemeral -c "cd ${COMPILE_REPO_LOCAL} && go test ./..."
```

### fmt
```sh
devenv ${DEVENV_FROM} ephemeral -c "cd ${COMPILE_REPO_LOCAL} && go fmt ./..."
```

### tidy
```sh
devenv ${DEVENV_FROM} ephemeral -c "cd ${COMPILE_REPO_LOCAL} && go mod tidy"
```
