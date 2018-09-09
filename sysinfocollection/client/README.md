工具作用：
   收集指定数据客户端
   默认收集数据是：磁盘分区占用百分比，cpu负载，内存使用率 

安装使用：
   解压后即可使用
   解压：
      $ tar zxvf mokylin_client.tar.gz

   使用方法：
      $ python mokylin_client_daemon.py (start|stop)

关于扩展
   在moudles目录下，可以自行添加想要的功能，并编辑mokylin_client.py，加入到返回数据中即可。


注意事项：
   1.配置文件
     配置说明，见配置文件注释。
   2.日志目录
     当前目录logs，日志名称默认为mokylin_client.log。
