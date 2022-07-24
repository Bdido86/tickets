# For developers

## tg bot


## gRPC - настройка, proto файлы 

### Генерация protoBuf
See manual [docs](https://github.com/grpc-ecosystem/grpc-gateway)  
1) Установить зависимости на компьютер(сервер)  
    >go install \  
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \  
    google.golang.org/protobuf/cmd/protoc-gen-go \  
    google.golang.org/grpc/cmd/protoc-gen-go-grpc>
3) Выполнить команду  `make .dev-tools` - сгенерирует в папке /bin файлы для proto buf
4) Выполнить команду  `make .dev-buf` - запустит генерацию proto файлов

Как и в какую папку кладуться файлы - можно ознакомистья в файле `buf.gen.yaml`



