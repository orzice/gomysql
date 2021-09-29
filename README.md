<h1 align="center">Mysql Query Constructor</h1>

<p align="center">
Mysql query builder similar to laravel
</p>

## Install
```
1. download MysqlConstructor.go
2. Revise package
```

## Use

### Init
```
DB := extend.DB_Mysql{} //init
DB.Construct()
```
### Get
```
DB.Table("Demo").Where("state","=","1").WhereOr("stete","=","2").Order("id","desc").Order("oid","").Limit(1,1).Get()

【mysql】
SELECT * FROM `Demo` WHERE  `state` = '1' OR `state` = '2'   ORDER BY `id` DESC,`oid`  LIMIT 1,1
```
### Insert
```
data := map[string]string{
    "name" : "Orzice",
    "github" : "https://github.com/orzice",
}
DB.Table("Demo").Insert(data)

【mysql】
INSERT INTO `Demo` (name , github) VALUES ('Orzice' , 'https://github.com/orzice')
```
### Update
```
data := map[string]string{
    "name" : "Orzice",
    "github" : "https://github.com/orzice",
}
DB.Table("Demo").Where("id","=","1").Update(data)

【mysql】
UPDATE `Demo` SET  name = 'Orzice' , github = 'https://github.com/orzice'  WHERE  `id` = '1'
```
### Delete
```
DB.Table("Demo").Where("id","=","1").Delete()
	
【mysql】
DELETE FROM `Demo` WHERE  `id` = '1'
```
### Increment
```
DB.Table("Demo").Where("id","=","2").Set("weight","weight+1").Set("time","time+1").Update(map[string]string{})
	
【mysql】
UPDATE `Demo` SET weight=weight+1,time=time+1   WHERE  `id` = '2'
```
### Close
```
DB.Close()
```


## Quote

https://github.com/go-sql-driver/mysql