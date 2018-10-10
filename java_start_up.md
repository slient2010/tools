# 我本地有一个 data-service.jar


## 具体操作

### 1. 编写启动脚本  data-service-start

```bash
[root@iz2ze0fq2isg8vphkpos5sz shell]# more  data-service-start
#!/bin/sh

export JAVA_HOME=/usr/local/jdk1.8.0_131
export PATH=$JAVA_HOME/bin:$PATH

java -jar /data/imgcloud/data-service.jar > /data/logs/data-service.log &
echo $! > /var/run/data-service.pid
```

### 2. 编写停止脚本

```bash
[root@iz2ze0fq2isg8vphkpos5sz shell]# more data-service-stop 
#!/bin/sh
PID=$(cat /var/run/data-service.pid)
kill -9 $PID
```
### 3. 在/usr/lib/systemd/system 下 编写 data-service.service 脚本

```bash
[root@iz2ze0fq2isg8vphkpos5sz shell]# cd /usr/lib/systemd/system
[root@iz2ze0fq2isg8vphkpos5sz system]# more data-service.service 
[Unit]
Description=data-service for mongodb
After=syslog.target network.target remote-fs.target nss-lookup.target
 
[Service]
Type=forking
ExecStart=/data/shell/data-service-start
ExecStop=/data/shell/data-service-stop
PrivateTmp=true
 
[Install]
WantedBy=multi-user.target
```

### 4. 相关命令

```
  systemctl  enable   data-service   #  开机自启动
  systemctl  stop  data-service  # 停止
  system  start data-service  # 启动
```

## 参考资料

>Ref

---------------------

原文：https://blog.csdn.net/stonexmx/article/details/73541508?utm_source=copy 

