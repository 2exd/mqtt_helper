使用 mosquitto 作为消息中转服务

## server

修改 `./temp.txt` 发送文本给 client 端剪切板

- 启动

```shell
./helper.exe server
```

## client

快捷键 `1+2+3` 截图发送给 server

- 启动

```shell
./helper.exe client
```

## 配置文件

一般不用修改。

也可以使用自己搭建的 mosquitto 服务器。



## TODO

- 连接管理功能

- server 私发给 client