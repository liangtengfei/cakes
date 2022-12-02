# 读写分离多数据库支持
### DBResolver 为 GORM 提供了多个数据库支持，支持以下功能：

1. 支持多个 sources、replicas
1. 读写分离
1. 根据工作表、struct 自动切换连接
1. 手动切换连接
1. Sources/Replicas 负载均衡
1. 适用于原生 SQL
1. 事务

### 配置方式
```
import (
  "gorm.io/gorm"
  "gorm.io/plugin/dbresolver"
  "gorm.io/driver/mysql"
)

db, err := gorm.Open(mysql.Open("db1_dsn"), &gorm.Config{})

db.Use(dbresolver.Register(dbresolver.Config{
  // `db2` 作为 sources，`db3`、`db4` 作为 replicas
  Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
  Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
  // sources/replicas 负载均衡策略
  Policy: dbresolver.RandomPolicy{},
}).Register(dbresolver.Config{
  // `db1` 作为 sources（DB 的默认连接），对于 `User`、`Address` 使用 `db5` 作为 replicas
  Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
}, &User{}, &Address{}).Register(dbresolver.Config{
  // `db6`、`db7` 作为 sources，对于 `orders`、`Product` 使用 `db8` 作为 replicas
  Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
  Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
}, "orders", &Product{}, "secondary"))
```
> 如果没有配置sources，第一个默认就是source
> 插件地址：https://github.com/go-gorm/dbresolver

# 分表支持
### Sharding 是一个高性能的 Gorm 分表中间件。它基于 Conn 层做 SQL 拦截、AST 解析、分表路由、自增主键填充，带来的额外开销极小。对开发者友好、透明，使用上与普通 SQL、Gorm 查询无差别，只需要额外注意一下分表键条，就能提供高性能的数据库访问。支持一下功能：
1. 非侵入式设计， 加载插件，指定配置，既可实现分表。
1. 轻快， 非基于网络层的中间件，像 Go 一样快
1. 支持多种数据库。 PostgreSQL 已通过测试，MySQL 和 SQLite 也在路上。
1. 多种主键生成方式支持（Snowflake, PostgreSQL Sequence, 以及自定义支持）Snowflake 支持从主键中确定分表键。

### 配置方式
```
import (
    "fmt"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/sharding"
)

dsn := "postgres://localhost:5432/sharding-db?sslmode=disable"
db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))

db.Use(sharding.Register(sharding.Config{
    ShardingKey:         "user_id",
    NumberOfShards:      64,
    PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders").Register(sharding.Config{
    ShardingKey:         "user_id",
    NumberOfShards:      256,
    PrimaryKeyGenerator: sharding.PKSnowflake,
    // This case for show up give notifications, audit_logs table use same sharding rule.
}, Notification{}, AuditLog{}))
```
> 数据表需要提前创建好。插件并不会自己创建表。
> 插件地址：https://github.com/longbridgeapp/gorm-sharding