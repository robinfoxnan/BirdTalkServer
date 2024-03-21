package core

import (
	"birdtalk/server/pbmodel"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"log"
)

// Encoder 接口定义了编码器的方法
type Encoder interface {
	Encode(data *pbmodel.Msg) ([]byte, error) // 编码函数，将数据编码为字节流
	Decode(data []byte) (*pbmodel.Msg, error) // 解码函数，将字节流解码为指定类型的数据
}

// JSONEncoder 实现了 Encoder 接口，使用 JSON 格式进行编码和解码
type JSONEncoder struct{}
type BinEncoder struct{}

// Encode 使用 JSON 编码数据
func (e *JSONEncoder) Encode(data *pbmodel.Msg) ([]byte, error) {
	// 将数据编码为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// Decode 使用 JSON 解码数据
func (e *JSONEncoder) Decode(data []byte) (*pbmodel.Msg, error) {
	// 将 JSON 数据解码为指定类型的数据
	msg := pbmodel.Msg{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// 解码函数
func (*BinEncoder) DecodeMsg(data []byte) (*pbmodel.Msg, error) {
	msg := pbmodel.Msg{}
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Fatalf("protobuf unmarshal error: %v", err)
	}
	return &msg, err
}

func (*BinEncoder) EncodeMsg(msg *pbmodel.Msg) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalf("protobuf marshal error: %v", err)
	}

	return data, err
}
