**发现一个非常好用的工具**

举例：搭建主从复制的PostgreSQL. 共三步。

1. 下载Docker
下载地址：https://www.docker.com/

2. 下载compose文件
```
$ curl -sSL https://raw.githubusercontent.com/bitnami/containers/main/bitnami/postgresql/docker-compose-replication.yml > docker-compose-replication.yml
```
3. 运行compose文件
```
docker-compose -f docker-compose-replication.yml up -d
```

完成以上步骤就可以直接从桌面访问了，搭建开发环境非常方便。
> 别忘了映射软件的端口。比如：ports: - '5432:5432'

更详细的介绍可以查看官网：https://bitnami.com [Loved by Developers.
Trusted by Ops](https://bitnami.com/)