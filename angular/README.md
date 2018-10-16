# angular cli 

创建项目

```
ng new projectname
# --flat
cd projenctname/src
ng generate module app/core/core --flat --spec=false
ng generate guard --flat app/core/module-import-guard --spec=false
ng generate module app/layout/layout --flat --spec=false
ng generate module app/shared/shared --flat --spec=false
ng generate module app/routes/routes --routing --flat --spec=false
```
or 
```
cd src
ng generate module app/core --spec=false
ng generate guard --flat app/core/module-import-guard --spec=false
ng generate module app/layout --spec=false
ng generate module app/shared --spec=false
ng generate module app/routes --routing --spec=false
```
