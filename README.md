# RedisHoneyPot
***
RedisHoneyPot是一款支持Redis协议的高交互式蜜罐系统。采用Golang语言开发。

* 已实现模拟命令
1. ping
2. info
3. set
4. get
5. del
6. exists
7. keys
8. flushall
9. flushdb
10. save
11. select
12. dbsize
13. config
14. slaveof

* 使用方法
``` nohup ./RedisHoneyPot -addr 0.0.0.0:6379 -proto tcp -num 1 > redis.log &```
  * 默认监听0.0.0.0的6379端口，使用tcp协议
