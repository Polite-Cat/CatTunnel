# PoliteCat

使用TUN转发流量的高效VPN，使用Go语言开发。

![image](https://img.shields.io/badge/License-Apache2.0-orange)

# 特性
* websocket通信
* 支持服务端流量记录
* 支持流量加密
* 持久更新

# 用法

```
Usage of politeCat:
  -S    server mode
  -f string
        mixin function xor/none (default "xor")
  -k string
        key (default "fuck_key")
  -path string
        ws path (default "/freedom")
  -s string
        server address (default ":3001")
  -t int
        dial timeout in seconds (default 30)
```

# 使用
## 基础使用
### 编译
```
go build main.go
```

### 客户端（转发所有流量）
```
sudo ./politeCat -s your_serverIP:3001
```


### 服务端
```
sudo ./politeCat -S
```

### ***在Linux上设置服务端***
***需执行如下代码：***
```
  echo 1 > /proc/sys/net/ipv4/ip_forward
  sysctl -p
  # 设置NAT转发流量
  iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
  iptables -t nat -A POSTROUTING -o tun0 -j MASQUERADE
  iptables -A INPUT -i eth0 -m state --state RELATED,ESTABLISHED -j ACCEPT
  iptables -A INPUT -i tun0 -m state --state RELATED,ESTABLISHED -j ACCEPT
  iptables -A FORWARD -j ACCEPT
  sysctl -p
```

### 移动端
**待实现**

# TODO
1. 支持Android
2. 支持更多协议
3. 动态分配Client的ip

## 进阶使用（<u>完善中</u>）

### 查看服务端流量使用情况

http://your_serverIP:3001/stats

### 查看服务端眼中你的ip

http://your_serverIP:3001/ip

### 查看服务端的ip列表

http://your_serverIP:3001/register/list/ip

# 感谢

感谢 [net-byte/vtun](https://github.com/net-byte/vtun/) 提供的思路。

