# angular cli 使用总结

## 创建项目
```bash
ng new cmdb-angular --spec=false
```

## 安装组件

```
# --flat
cd cmdb-angular/src
ng generate module app/core/core --flat --spec=false
ng generate guard --flat app/core/module-import-guard --spec=false
ng generate module app/layout/layout --flat --spec=false
ng generate module app/shared/shared --flat --spec=false
ng generate module app/routes/routes --routing --flat --spec=false
```

or 

## 安装组件
```bash
# without --flat
cd cmdb-angular/src
ng generate module app/core --spec=false
ng generate guard --flat app/core/module-import-guard --spec=false
ng generate module app/layout --spec=false
ng generate module app/shared --spec=false
ng generate module app/routes --routing --spec=false
```


## 增加模块
```bash
ng add ng-zorro-antd --theme --styleex=less
# 登录框架
cd app/routes
ng g ng-zorro-antd:form-normal-login passport --styleext=less --name=login --spec=false

```

## 增加拦截器

```bash
ng g class app/core/net/default.interceptor  --spec=false
```
