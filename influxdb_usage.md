# influxdb基本操作


## 名词解释

在具体的讲解influxdb的相关操作之前先说说influxdb的一些专有名词，这些名词代表什么。

influxDB名词
```
database：数据库；

measurement：数据库中的表；

points：表里面的一行数据。

influxDB中独有的一些概念

Point由时间戳（time）、数据（field）和标签（tags）组成。

time：每条数据记录的时间，也是数据库自动生成的主索引；

fields：各种记录的值；

tags：各种有索引的属性。

series: 所有在数据库中的数据，都需要通过图表来表示，series表示这个表里面的所有的数据可以在图标上画成几条线（注：线条的个数由tags排列组合计算出来）
```

举个简单的小例子：
有以下数据：
```
> select * from cpu_usage ;
name: cpu_usage
time                cpu       idle system user
----                ---       ---- ------ ----
1537239708000000000 cpu-total 10.1 53.3   46.6
1537239801000000000 cpu-total 10.1 53.3   46.6
1537239805000000000 cpu-total 10.1 53.3   46.6
1537239806000000000 cpu-total 10.1 53.3   46.6
1537239807000000000 cpu-total 10.1 53.3   46.6
1537239808000000000 cpu-total 10.1 53.3   46.6
1537239809000000000 cpu-total 10.1 53.3   46.6
1537239810000000000 cpu-total 10.1 53.3   46.6
1537239812000000000 cpu-total 10.1 53.3   46.6
1537239813000000000 cpu-total 10.1 53.3   46.6
1537239814000000000 cpu-total 10.1 53.3   46.6
1537239815000000000 cpu-total 10.1 53.3   46.6
1537239816000000000 cpu-total 10.1 53.3   46.6
1537239830000000000 cpu-total 10.1 53.3   46.6
1537239831000000000 cpu-total 10.1 53.3   46.6
1537239832000000000 cpu-total 10.1 53.3   46.6
1537239833000000000 cpu-total 10.1 53.3   46.6
1537239834000000000 cpu-total 10.1 53.3   46.6
1537239835000000000 cpu-total 10.1 53.3   46.6
1537239836000000000 cpu-total 10.1 53.3   46.6
1537239837000000000 cpu-total 10.1 53.3   46.6
1537239838000000000 cpu-total 10.1 53.3   46.6
1537239839000000000 cpu-total 10.1 53.3   46.6
1537239840000000000 cpu-total 10.1 53.3   46.6
1537239841000000000 cpu-total 10.1 53.3   46.6
1537239857000000000 cpu-total 10.1 53.3   46.6
1537239858000000000 cpu-total 10.1 53.3   46.6
1537239859000000000 cpu-total 10.1 53.3   46.6
```

它的series为：
```
> show series from cpu_usage ;
key
---
cpu_usage,cpu=cpu-total
> 

```


## influxDB基本操作

### 数据库与表的操作

可以直接在web管理页面做操作，当然也可以命令行。

```
#创建数据库
create database "db_name"

#显示所有的数据库
show databases

#删除数据库
drop database "db_name"

#使用数据库
use db_name

#显示该数据库中所有的表
show measurements

#创建表，直接在插入数据的时候指定表名
insert test,host=127.0.0.1,monitor_name=test count=1

#删除表
drop measurement "measurement_name"
```

#### 增

向数据库中插入数据。

通过命令行

```
 > use metrics
Using database metrics
> insert test,host=127.0.0.1,monitor_name=test count=1
```
这样，数据库插入数据成功。

通过http接口

```bash
$ curl -i -XPOST 'http://127.0.0.1:8086/write?db=metrics' --data-binary 'test,host=127.0.0.1,monitor_name=test count=1'
```


读者看到这里可能会观察到插入的数据的格式貌似比较奇怪，这是因为influxDB存储数据采用的是Line Protocol格式。那么何谓Line Protoco格式？

Line Protocol格式：写入数据库的Point的固定格式。

在上面的两种插入数据的方法中都有这样的一部分：
```
test,host=127.0.0.1,monitor_name=test count=1
```
其中：

```
test：表名；
host=127.0.0.1,monitor_name=test：tag；
count=1：field
```
相对此格式有详细的了解参见官方文档

#### 查

查询数据库中的数据。

通过命令行

```
 > use metrics
Using database metrics
> select * from test order by time desc
```
通过http接口
```bash
$ curl -G 'http://localhost:8086/query?pretty=true' --data-urlencode "db=metrics" --data-urlencode "q=select * from test order by time desc"
```


influxDB是支持类sql语句的，具体的查询语法都差不多，这里就不再做详细的赘述了。

### 数据保存策略（Retention Policies）

influxDB是没有提供直接删除数据记录的方法，但是提供数据保存策略，主要用于指定数据保留时间，超过指定时间，就删除这部分数据。

查看当前数据库Retention Policies

```
show retention policies on "db_name"
```

### 创建新的Retention Policies

```
create retention policy "rp_name" on "db_name" duration 3w replication 1 default
```

```
rp_name：策略名
db_name：具体的数据库名
3w：保存3周，3周之前的数据将被删除，influxdb具有各种事件参数，比如：h（小时），d（天），w（星期）
replication 1：副本个数，一般为1就可以了
default：设置为默认策略
```

### 修改Retention Policies

```
alter retention policy "rp_name" on "db_name" duration 30d default
```

### 删除Retention Policies

```
drop retention policy "rp_name"
```


### 连续查询（Continous Queries）

当数据超过保存策略里指定的时间之后就会被删除，但是这时候可能并不想数据被完全删掉，怎么办？

influxdb提供了联系查询，可以做数据统计采样。

查看数据库的Continous Queries

```
show continuous queries
```


### 创建新的Continous Queries

```
create continous query cq_name on db_name begin select sum(count) into new_table_name from table_name group by time(30m) end
```

```
cq_name：连续查询名字
db_name：数据库名字
sum(count)：计算总和
table_name：当前表名
new_table_name：存新的数据的表名
30m：时间间隔为30分钟
```

### 删除Continous Queries

```
drop continous query cp_name on db_name
```

### 用户管理

可以直接在web管理页面做操作，也可以命令行。

```
#显示用户
show users

#创建用户
create user "username" with password 'password'

#创建管理员权限用户
create user "username" with password 'password' with all privileges

#删除用户
drop user "username"
```
