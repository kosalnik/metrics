Для генерации кода из этого proto файла нужно установить утилиты:

```shell
sudo apt install protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Затем компилируем:
```shell
protoc --experimental_allow_proto3_optional \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  ./pkg/metrics/metrics.proto
```

Для использования того что скомпилировано нужно подключить нужные пакеты:
```shell
go get google.golang.org/protobuf
go get google.golang.org/grpc
```