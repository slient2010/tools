# 小技巧：ssh隧道+iptables 实现外部访问内部windows服务器的远程桌面

```bash
#!/bin/bash
num=`ps aux |grep "ServerAliveInterval=180 user@public.server"| grep 3389 |grep -v grep |wc -l`
[ ${num} -eq 0 ] && ssh -o ServerAliveInterval=180  user@public.server -p 22 -R 0.0.0.0:3389:inner.server:3389 -fN
n=`sudo iptables -L -n -t nat |grep 3389 | grep -v grep |wc -l`
# 在内网任意一台服务器（与目标主机能通讯的主机）执行请求转发
[ ${num} -eq 0 ] && sudo iptables -t nat -A PREROUTING -d inner.server -p tcp --dport 3389 -j DNAT --to inner.server:3389
```

public.server => 公网服务器IP
inner.server => 内网服务器IP
