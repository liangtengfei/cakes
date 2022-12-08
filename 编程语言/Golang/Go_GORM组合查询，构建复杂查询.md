# 组合
使用 Group 条件可以更轻松的编写复杂 SQL
```GO
db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium"))
```
形成的SQL语句：
```SQL
WHERE pizza = "pepperoni" AND (size = "small" OR size = "medium")
```

# 子查询
子查询可以嵌套在查询中，GORM 允许在使用 *gorm.DB 对象作为参数时生成子查询
```
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
```
形成的SQL语句：
```SQL
SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");
```

# From 子查询
GORM 允许您在 Table 方法中通过 FROM 子句使用子查询，例如：
```
subQuery1 := db.Model(&User{}).Select("name")
subQuery2 := db.Model(&Pet{}).Select("name")
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
```
形成的SQL语句：
```SQL
SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
```

> 记录一下：引用自[GORM](https://gorm.io/zh_CN/docs/advanced_query.html)