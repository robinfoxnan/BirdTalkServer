.\protoc.exe --proto_path=. user.proto msg.proto --js_out=import_style=typescript:../server/js/model
protoc -I. --go_out=../ --go-grpc_out=../ user.proto msg.proto