# Openops
## 开发环境
```
golang, python2.6(elasticsearch api)
golang环境搭建可以参考网上

此处GOPATH=/data/
项目目录：/data/src/openops
```
## 相关依赖下载
```
> cd $GOPATH
# 数据库接口
> go get github.com/ziutek/mymysql/mysql
> go get github.com/ziutek/mymysql/native

# session 
> go get github.com/gorilla/sessions
> go get github.com/gorilla/context
```

## 项目安装
```
1.安装go-bindata工具，linux下，最新的源仓库应该可以搜的到。
go-bindata
Install go-bindata
> apt-get install go-bindata

2.安装完成后，如下打包项目，生成的二进制文件可以放到平台直接使用。
> cd front
> go-bindata -pkg=front -nocompress=true -debug=true html/...
> go-bindata -pkg=front -nocompress=true html/...
> go install openops
```
## 数据库
```
# 数据库在源码的/data目录下，直接导入数据库中就可以
# 添加用户信息到user中就可以登录
mysql> insert into user value('', 'xxx@xx.cn', 'password', '孟洋', '17075000000', '2016-12-07 11:03:00', '','', '');

注：密码明文
登录后台直接使用该密码
```

## 项目执行
```
> cd $GOPATH/bin
> ./openops 
```

## 预览
![界面预览](https://raw.githubusercontent.com/slient2010/openops/master/demo_images_1.png)



参考地址:

[http://studygolang.com/articles/884](http://studygolang.com/articles/884)

[https://github.com/josephspurrier/gowebapp](https://github.com/josephspurrier/gowebapp) 
