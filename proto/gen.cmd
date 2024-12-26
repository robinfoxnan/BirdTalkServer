.\protoc.exe --proto_path=. user.proto msg.proto --js_out=import_style=typescript:../server/js/model
protoc -I. --go_out=../ --go-grpc_out=../ user.proto msg.proto
./protoc --java_out=./java user.proto msg.proto
protoc --python_out=./py user.proto msg.proto
./protoc.exe --kotlin_out=./java user.proto msg.proto