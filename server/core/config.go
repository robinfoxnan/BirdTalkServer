package core

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var Config *LocalConfig = nil

// 手动调用这个"config.yaml"
func InitConfig(filename string) error {
	Globals.Logger.Info("load config", zap.String("name", filename))

	var err error
	Config, err = LoadConfig(filename)
	return err
}

type RedisConf struct {
	RedisHost string `yaml:"redis_host" `
	RedisPwd  string `yaml:"redis_pwd"`
}

type MongoConf struct {
	MongoHost string `yaml:"mongo_host" `
	DbName    string `yaml:"db_name"`
}

type ScyllaConf struct {
	Host string `yaml:"scylla_host"`
	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`
}

type EmailConf struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`
	Tls  bool   `yaml:"tls"`
}

type ServerConf struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	HostIndex int    `yaml:"host_index"`
	HostName  string `yaml:"host_name"`
	MsgQueLen int    `yaml:"group_msg_queue_len"`

	Workers    int       `yaml:"workers"`
	Schema     string    `yaml:"schema"`
	CertFile   string    `yaml:"cert"`
	KeyFile    string    `yaml:"key"`
	FriendMode bool      `yaml:"friend_making"`
	Email      EmailConf `yaml:"email"`
}

type LocalConfig struct {
	Redis    RedisConf   `yaml:"redis"`
	Server   ServerConf  `yaml:"server"`
	MongoDb  MongoConf   `yaml:"mongoDb"`
	ScyllaDb ScyllaConf  `yaml:"scyllaDb"`
	Email    EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	SMTPAddr              string `yaml:"smtp_addr"`
	SMTPPort              string `yaml:"smtp_port"`
	SMTPHeloHost          string `yaml:"smtp_helo_host"`
	UserName              string `yaml:"user_name"`
	UserPwd               string `yaml:"user_pwd"`
	TLSInsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
	AuthType              string `yaml:"auth_type"`
}

// "config.yaml"
func LoadConfig(fileName string) (*LocalConfig, error) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	} // 将读取的yaml文件解析为响应的 struct
	config := LocalConfig{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	if config.Server.MsgQueLen < 100 {
		config.Server.MsgQueLen = 100
	}
	if config.Server.Workers < 1 {
		config.Server.Workers = 1
	}
	return &config, err
}

func SaveConfig(conf *LocalConfig) bool {
	data, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(string(data))
	err = os.WriteFile("./config.yaml", data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return true
}
