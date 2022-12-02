> Redis 在 Windows 上不受官方支持。但是，可以按照以下说明在 Windows 上安装 Redis 进行开发。借用WSL。

要在 Windows 上安装 Redis，首先需要启用 WSL2（适用于 Linux 的 Windows 子系统）。 WSL2 允许在 Windows 上本地运行 Linux 二进制文件。要使此方法起作用，需要运行 Windows 10 版本 2004 及更高版本或 Windows 11。
# 安装Redis
在 Windows 上运行 Ubuntu 后，可以按照在 Ubuntu/Debian 上安装中详述的步骤从官方 packages.redis.io APT 存储库安装 Redis 的最新稳定版本。将存储库添加到 apt 索引，更新它，然后安装：
```shell
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

sudo apt-get update
sudo apt-get install redis
```
最后，像这样启动 Redis 服务器：
```
sudo service redis-server start
```
# 连接Redis
可以通过连接 Redis CLI 来测试 Redis 服务器是否正在运行：
```
redis-cli 
127.0.0.1:6379> ping
PONG
```