# PoliteCat

使用TUN转发流量的高效VPN，使用Go语言开发。

![image](https://img.shields.io/badge/License-Apache2.0-orange)

# 特性
* 支持tcp
* 支持服务端流量记录
* 支持流量加密
* 持久更新

# 用法

```
Usage of ./politeCat:
  -S	server mode
  -c string
    	tun interface cidr (default "172.16.0.10/24")
  -d string
    	dns address (default "8.8.8.8:53")
  -enc
    	enable data encry
  -g	client global mode
  -k string
    	key (default "123456")
  -l string
    	local address (default ":3000")
  -mtu int
    	tun mtu (default 1500)
  -s string
    	server address (default ":3001")
  -t int
    	dial timeout in seconds (default 30)
```

## 编译

```
go build main.go
```

## 客户端

```
sudo ./politeCat -l=:3000 -s=server-addr:3001 -c=172.16.0.10/24 -k=123456 -enc
```

## 全局模式客户端（转发所有流量）

```
sudo ./politeCat -l=:3000 -s=server-addr:3001 -c=172.16.0.10/24 -k=123456 -g -enc
```

## 服务端

```
sudo ./politeCat -S -l=:3001 -c=172.16.0.1/24 -k=123456 -enc
```

## 在Linux上设置服务端

需执行如下代码：

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

## 移动端

**待实现**

## TODO
1. 支持Windows
2. 支持Android
3. 支持websocket
4. 部署完整的前后端

## 感谢

感谢 [net-byte/vtun](https://github.com/net-byte/vtun/) 提供的思路。

