# CatTunnel (原PoliteCat)

基于TUN网卡技术的C-S通信方案，使用Go语言开发。

![image](https://img.shields.io/badge/License-Apache2.0-orange)

# 应用场景
- 远程访问：服务端可通过客户端的TUN网卡ip访问到客户端的接口。
- 安全通信：客户端与服务端进行双向通信时，传输流量是加密状态，需要认证。
- 虚拟局域网（VLAN）扩展：使得不同物理位置的设备可以通过TUN网卡连接到同一虚拟网络中。

# 特性
* websocket通信
* 服务端流量记录
* 流量加密
* 持久更新
* 全平台支持

# 用法

```
Usage of catTunnel:
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

### 客户端（转发所有流量；在windows上需使用管理员运行）
```
sudo ./catTunnel -s your_serverIP:3001
```


### 服务端
```
sudo ./catTunnel -S
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

