# mysql中的外键foreign key

## 一、外键定义及使用

如果一张表中有一个非主键的字段指向了别一张表中的主键，就将该字段叫做外键。

　　一张表中可以有多个外键。

　　外键的默认作用有两点：

　　1.对子表(外键所在的表)的作用：子表在进行写操作的时候，如果外键字段在父表中找不到对应的匹配，操作就会失败。

　　2.对父表的作用：对父表的主键字段进行删和改时，如果对应的主键在子表中被引用，操作就会失败。

　　外键的定制作用----三种约束模式：

　　　　district：严格模式(默认), 父表不能删除或更新一个被子表引用的记录。

　　　　cascade：级联模式, 父表操作后，子表关联的数据也跟着一起操作。

　　　　set null：置空模式,前提外键字段允许为NLL,  父表操作后，子表对应的字段被置空。

　　使用外键的前提：

　　1. 表储存引擎必须是innodb，否则创建的外键无约束效果。

　　2. 外键的列类型必须与父表的主键类型完全一致。

　　3. 外键的名字不能重复。

　　4. 已经存在数据的字段被设为外键时，必须保证字段中的数据与父表的主键数据对应起来。

## 二、新增外键

　　1. 在创建时增加

　　　create table my_tab1(

　　　id int primary key auto_increment,

　　　name varchar(10) not null,

　　　class int,

　　　foreign key(class)　references my_tab2(主键字段名);

　　　)charset utf8;

　　2. 在创建好的表中增加

　　　alter table my_tab1 add [constraint 外键名] foreign key(外键字段名) references mytab2(主键字段名);

## 三、删除外键

　　alter table my_tab drop foreign key 外键名字;


## 四、参考资料
>[https://www.cnblogs.com/pengyin/p/6375860.html](https://www.cnblogs.com/pengyin/p/6375860.html)

