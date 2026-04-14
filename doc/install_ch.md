# 1. Redis安装流程

---

## 1.1 先判断你的 Linux 是什么系统
执行下面命令看结果：
```bash
cat /etc/os-release
```

- 显示 `CentOS` / `Red Hat` → 用 **CentOS 教程**
- 显示 `Ubuntu` / `Debian` → 用 **Ubuntu 教程**

---

## 1.2 Ubuntu / Debian 安装 Redis（apt 方式）
### 1.2.1. 更新源
```bash
sudo apt update -y
sudo apt upgrade -y
```

### 1.2.2. 直接安装
```bash
sudo apt install redis-server -y
```

### 1.2.3. 启动并设置开机自启
```bash
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

### 1.2.4. 查看状态
```bash
sudo systemctl status redis-server
```

### 1.2.5. 测试连接
```bash
redis-cli
127.0.0.1:6379> ping
```
返回 `PONG` 就是成功。

---

## 1.3 CentOS 7 / 8 / 9 安装 Redis（yum / dnf）
### 1.3.1. 安装 EPEL 源（必须）
```bash
sudo yum install epel-release -y
```

### 1.3.2. 安装 Redis
```bash
sudo yum install redis -y
```

### 1.3.3. 启动 + 开机自启
```bash
sudo systemctl start redis
sudo systemctl enable redis
```

### 1.3.4. 查看状态
```bash
sudo systemctl status redis
```

### 1.3.5. 测试
```bash
redis-cli
127.0.0.1:6379> ping
```

---

## 1.4 最关键：允许远程访问（生产/局域网常用）
**Ubuntu 和 CentOS 改法一样**

### 1.4.1. 打开配置文件

```bash
sudo vim /etc/redis/redis.conf
```

### 1.4.2. 修改两处
```
bind 0.0.0.0      # 允许所有IP访问
protected-mode no # 关闭保护模式
```

### 1.4.3. 重启 Redis
```bash
# Ubuntu
sudo systemctl restart redis-server

# CentOS
sudo systemctl restart redis
```

### 1.4.4. 防火墙放行 6379
#### Ubuntu
```bash
sudo ufw allow 6379
```

#### CentOS
```bash
sudo firewall-cmd --add-port=6379/tcp --permanent
sudo firewall-cmd --reload
```

---

## 1.5 设置密码（强烈建议）
编辑配置：
```bash
sudo vim /etc/redis/redis.conf
```

找到：
```
# requirepass foobared
```
改成：
```
requirepass 你的密码
```

重启：
```bash
# Ubuntu
sudo systemctl restart redis-server

# CentOS
sudo systemctl restart redis
```

连接时：
```bash
redis-cli -a 密码
```

---

## 1.6 一句话区分 CentOS 和 Ubuntu 命令
| 动作 | Ubuntu | CentOS |
|------|--------|--------|
| 安装 | apt install redis-server | yum install redis |
| 服务名 | redis-server | redis |
| 重启 | systemctl restart redis-server | systemctl restart redis |
| 防火墙 | ufw | firewall-cmd |

---

# 2. mongoDb安装
我给你一套**CentOS / Ubuntu 通用、最标准、直接复制运行**的 **MongoDB 安装教程**，和刚才 Redis 格式完全一样，**不踩坑、一步到位**。

## 2.1 先判断系统
```bash
cat /etc/os-release
```
看到 **Ubuntu** → 用 Ubuntu 教程
看到 **CentOS** → 用 CentOS 教程

---

## 2.2 Ubuntu 安装 MongoDB（官方源）
### 2.2.1. 导入公钥
```bash
sudo apt-get install gnupg curl -y
curl -fsSL https://pgp.mongodb.com/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor
```

### 2.2.2. 添加源
```bash
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
```

### 2.2.3. 安装
```bash
sudo apt update -y
sudo apt install mongodb-org -y
```

### 2.2.4. 启动 & 开机自启
```bash
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 2.2.5. 查看状态
```bash
sudo systemctl status mongod
```

### 2.2.6. 测试
```bash
mongosh
```

---

## 2.3 CentOS 安装 MongoDB
### 2.3.1. 添加 yum 源
```bash
sudo vim /etc/yum.repos.d/mongodb-org-7.0.repo
```

写入内容：
```
[mongodb-org-7.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/7.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-7.0.asc
```

### 2.3.2. 安装
```bash
sudo yum install mongodb-org -y
```

### 2.3.3. 启动 & 开机自启
```bash
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 2.3.4. 查看状态
```bash
sudo systemctl status mongod
```

### 2.3.5. 测试
```bash
mongosh
```

---

## 2.4 允许远程访问（两个系统一样）
### 2.4.1. 编辑配置
```bash
sudo vim /etc/mongod.conf
```

找到：
```
net:
  port: 27017
  bindIp: 127.0.0.1
```

改为：
```
net:
  port: 27017
  bindIp: 0.0.0.0
```

### 2.4.2. 重启
```bash
sudo systemctl restart mongod
```

### 2.4.3. 防火墙放行
**Ubuntu**
```bash
sudo ufw allow 27017
```

**CentOS**
```bash
sudo firewall-cmd --add-port=27017/tcp --permanent
sudo firewall-cmd --reload
```

---

## 2.5 设置用户名密码（重要）
进入 mongosh：
```bash
mongosh
```

切换库：
```
use admin
```

创建管理员：
```
db.createUser({
  user: "admin",
  pwd: "你的密码",
  roles: ["root"]
})
```

开启验证：
编辑 `/etc/mongod.conf`
```
security:
  authorization: enabled
```

重启：
```bash
sudo systemctl restart mongod
```

---

## 2.6 一句话区分 CentOS / Ubuntu MongoDB
| 操作 | Ubuntu | CentOS |
|------|--------|--------|
| 源文件 | /etc/apt/sources.list.d/xxx | /etc/yum.repos.d/xxx |
| 安装命令 | apt install mongodb-org | yum install mongodb-org |
| 服务名 | mongod | mongod |
| 配置文件 | /etc/mongod.conf | /etc/mongod.conf |
| 端口 | 27017 | 27017 |

# 3. 安装scyllaDb单机版

## 3.1 先判断系统
```bash
cat /etc/os-release
```
- 显示 **CentOS/RHEL** → 用 CentOS 教程
- 显示 **Ubuntu/Debian** → 用 Ubuntu 教程

---

## 3.2 Ubuntu 安装 ScyllaDB（20.04+/22.04+）
### 3.2.1. 安装依赖 & 导入 GPG 密钥
```bash
sudo apt update -y
sudo apt install -y gnupg2 curl
sudo mkdir -p /etc/apt/keyrings
sudo gpg --homedir /tmp --no-default-keyring --keyring /etc/apt/keyrings/scylladb.gpg \
  --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys A43E06657BAC99E3 
```

### 3.2.2. 添加官方源（2025.4 稳定版）
```bash
sudo wget -O /etc/apt/sources.list.d/scylla.list \
  http://downloads.scylladb.com/deb/debian/scylla-2025.4.list 
```

### 3.2.3. 安装 ScyllaDB
```bash
sudo apt update -y
sudo apt install -y scylla 
```

### 3.2.4. 执行系统优化（必须！单机必跑）
```bash
sudo scylla_setup
```
- 全程按 **Enter** 选默认（单机用默认配置即可）
- 脚本自动优化：CPU、IO、内存、网络、磁盘挂载等

### 3.2.5. 启动 & 开机自启
```bash
sudo systemctl start scylla-server
sudo systemctl enable scylla-server
```

### 3.2.6. 查看状态
```bash
sudo systemctl status scylla-server
```

### 3.2.7. 测试连接（CQL Shell，兼容 Cassandra）
```bash
cqlsh
```
出现 `Connected to` 即成功。

---

## 3.3 CentOS 7 / 8 / 9 安装 ScyllaDB
### 3.3.1. 安装依赖
```bash
sudo yum install -y curl gpg
```

### 3.3.2. 添加官方 YUM 源（2025.4 稳定版）
```bash
sudo curl -o /etc/yum.repos.d/scylla.repo \
  https://downloads.scylladb.com/rpm/centos/scylla-2025.4.repo 
```

### 3.3.3. 安装
```bash
sudo yum install -y scylla 
```

### 3.3.4. 系统优化（必须）
```bash
sudo scylla_setup
```
- 全程 **Enter** 默认即可

### 3.3.5. 启动 & 自启
```bash
sudo systemctl start scylla-server
sudo systemctl enable scylla-server
```

### 3.3.6. 状态 & 测试
```bash
sudo systemctl status scylla-server
cqlsh
```

---

## 3.4 允许远程访问（单机常用）
### 3.4.1. 编辑主配置
```bash
sudo vim /etc/scylla/scylla.yaml
```

修改：
```yaml
listen_address: 0.0.0.0
rpc_address: 0.0.0.0
# 单机可注释 seeds（默认 127.0.0.1）
# seeds: "127.0.0.1"
```

### 3.4.2. 防火墙放行端口（9042 CQL / 9160 Thrift / 10000 JMX）
#### Ubuntu
```bash
sudo ufw allow 9042
sudo ufw allow 9160
sudo ufw allow 10000
```

#### CentOS
```bash
sudo firewall-cmd --add-port=9042/tcp --permanent
sudo firewall-cmd --add-port=9160/tcp --permanent
sudo firewall-cmd --add-port=10000/tcp --permanent
sudo firewall-cmd --reload
```

### 3.4.3. 重启服务
```bash
sudo systemctl restart scylla-server
```

---

## 3.5 设置用户名密码（安全必做）
### 3.5.1. 进入 cqlsh
```bash
cqlsh
```

### 3.5.2. 创建超级管理员
```cql
CREATE ROLE admin WITH SUPERUSER = true AND LOGIN = true AND PASSWORD = '你的密码';
```

### 3.5.3. 开启验证
```bash
sudo vim /etc/scylla/scylla.yaml
```

修改：
```yaml
authenticator: PasswordAuthenticator
authorizer: CassandraAuthorizer
```

### 3.5.4. 重启
```bash
sudo systemctl restart scylla-server
```

### 3.5.5. 带密码登录
```bash
cqlsh -u admin -p 你的密码
```

---

## 3.6 CentOS / Ubuntu 关键区别
| 项目 | Ubuntu | CentOS |
|------|--------|--------|
| 源文件 | `/etc/apt/sources.list.d/scylla.list` | `/etc/yum.repos.d/scylla.repo` |
| 安装命令 | `apt install scylla` | `yum install scylla` |
| 服务名 | `scylla-server` | `scylla-server` |
| 配置文件 | `/etc/scylla/scylla.yaml` | `/etc/scylla/scylla.yaml` |
| 默认端口 | 9042 (CQL) | 9042 (CQL) |
| 优化脚本 | `scylla_setup` | `scylla_setup` |

---

## 3.7 常用命令
```bash
# 状态
sudo systemctl status scylla-server
# 重启
sudo systemctl restart scylla-server
# 停止
sudo systemctl stop scylla-server
# 查看节点状态
nodetool status
# 日志
sudo tail -f /var/log/scylla/scylla.log
```

# 4. BirdTalkServer 

## 4.1 编译

```bash
cd BirdTalkServer
go mod tidy
cd server
go build
```

## 4.2 配置

使用图形界面编辑`gedit config.yaml`

备注：

1) 各个数据库建议仅支持本地的127.0.0.1连接，这样更安全；可以设置口令也可以没有口令；
2)  server部分：cert和key是https使用的证书和私钥，file_base_path是各种图片与附件以及长语音保存的路径；
3) token_key暂时没有使用；后续对客户端授权使用；
4) email部分：是邮箱登陆的发送验证码所使用的第三方服务，这里根据各个邮箱的配置来设置。



```yaml
redis:
    redis_host: 127.0.0.1:6379
    redis_pwd:
mongoDb:
    #mongo_host: mongodb://admin:123456@127.0.0.1:27017
    mongo_host: mongodb://127.0.0.1:27017
    db_name: birdtalk
scyllaDb:
    scylla_host: 8.140.203.92:9042
    user: cassandra
    pwd: 123456
kafka:
    kafka_host:

server:
    host: 0.0.0.0
    port: 7817
    host_index: 1
    host_name: node1
    group_msg_queue_len: 100
    workers: 1
    schema: https
    cert: ./certs/cert.pem
    key: ./certs/key.pem
    friend_making: false
    cluster_mode: false
    avatar_font: ./ttf/SourceHanSans-VF.ttf
    file_base_path: e:/upload
    geolite2_path: ./resource/GeoLite2-City.mmdb
    log_level: debug
    token_key: 278b3e9f1c18f488314c7991e163c8bac880d5f09f4a15dc7589823dc6b43264


email:
    #tjj.31415
    smtp_addr: smtp.sina.com
    smtp_port: 25
    smtp_helo_host: sina.com
    user_name: birdchat@sina.com
    user_pwd: 1234567890
    tls_insecure_skip_verify: false
    auth_type: login
    workers: 3


```

