# 基于MQTT的聊天室

<u>特别提醒：该内容仅供个人娱乐所用，请勿将其用于任何违反法律的用途，一旦出现问题本人概不负责。</u>

基于 [MQTT](https://mqtt.org/) 协议，使用 [mosquitto](https://mosquitto.org/) 服务器提供消息中转服务。

## server

修改 `./temp.txt` 广播发送文本给 client 端剪切板

- 启动

```shell
./helper.exe server
```

## client

快捷键 `1+2+3` 屏幕截图发送给 server

快捷键 `q+w+e` 将剪切板内容发送给 server

- 启动

```shell
./helper.exe client
```

## 配置文件

一般不用修改。

测试用 mqtt 服务器地址 `tcp://124.221.165.250:1883`

也可以使用自己搭建的 mosquitto 服务器。

## 编译

```shell
$OS_ARCH="windows"
$env:GOOS=${OS_ARCH}
$env:GOARCH="amd64"

go build -ldflags "-s -w" -o ./helper.exe ./main.go
```

## docker安装mosquitto

- 准备配置文件

```shell
echo -e "listener 1883\npersistence true\nallow_anonymous true\npersistence_file mosquitto.db\npersistence_location /mosquitto/data/\nlog_dest file /mosquitto/log/mosquitto.log" > /root/docker_files/mosquitto/mosquitto.conf
```

- 创建目录

```shell
mkdir -p /usr/local/services/mosquitto/mosquitto-log
mkdir -p /usr/local/services/mosquitto/mosquitto-data
```

- 启动容器

```shell
docker run -it -p 1883:1883 -v /root/docker_files/mosquitto:/mosquitto/config -v /usr/local/services/mosquitto/mosquitto-data:/mosquitto/data -v /usr/local/services/mosquitto/mosquitto-log:/mosquitto/log eclipse-mosquitto
```



## TODO

- server 私信 client
