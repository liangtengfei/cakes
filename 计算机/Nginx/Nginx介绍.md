# Nginx 简介

Nginx 作为一款面向性能设计的HTTP服务器，相较于Apache、lighttpd具有占有内存少，稳定性高等优势。其流行度越来越高，应用也越来越广泛，常见的应用有：网页服务器、反向代理服务器以及电子邮件（IMAP/POP3）代理服务器，高并发大流量站点常用来做接入层的[负载均衡](https://cloud.tencent.com/product/clb?from=20065&from_column=20065)，还有非常常见的用法是作为日志采集服务器等。

Nginx 整体采用模块化设计，有丰富的模块库和第三方模块库，配置灵活。其中模块化设计是nginx的一大卖点，甚至http服务器核心功能也是一个模块。要注意的是：nginx的模块是静态的，添加和删除模块都要对nginx进行重新编译，这一点与Apache的动态模块完全不同。不过后来淘宝做了二次开发开源的 tengine 是支持 官方所有的 HTTP 模块动态加载而不必重新编译 Nginx，除非是第三方模块才需要重新编译。因此，在生产环境中，推荐用淘宝开源的 tengine，本文也以 tengine 作为示例。

虽然 Nginx 有如此强大的性能以及众多的三方模块支持，但每次重新编译以及寻找三方模块对生产环境来说还是不可接受的，幸运的是，Nginx 它是支持客户自己 Lua 脚本编程扩展相应的功能的，而且可以热加载，这就给生产环境带来了无限可能。比如我现在想要直接用Nginx + redis 做反爬虫和频率限制，Nginx + Kafka 做日志的实时流处理等等。

注：lvs 和 nginx 的负载均衡区别：

LVS：Linux Virtual Server，基于IP的负载均衡和反向代理技术，所以它几乎可以对所有应用做负载均衡，包括http、数据库、在线聊天室等等，LVS工作在4层，在Linux内核中作四层交换，只花128个字节记录一个连接信息，不涉及到文件句柄操作，故没有65535最大文件句柄数的限制。LVS性能很高，可以支持100～400万条并发连接。抗负载能力强、是工作在网络4层之上仅作分发之用，**没有流量的产生**，这个特点也决定了它在负载均衡软件里的性能最强的，对内存和cpu、IO资源消耗比较低。 

Nginx：基于HTTP的负载均衡和反向代理服务器，Nginx工作在网络的7层，所以它可以针对http应用本身来做分流策略，比如针对域名、URL、目录结构等，相比之下LVS并不具备这样的功能，能够很好地支持[虚拟主机](https://cloud.tencent.com/product/lighthouse?from=20065&from_column=20065)，可配置性很强，大约能支持3～5万条并发连接。



# 关于 nginx 正则说明

（1）location 匹配语法规则

Nginx location 的正则匹配语法与优先级容易让新同学迷惑。

~    #波浪线表示执行一个正则匹配，区分大小写

~*   #表示执行一个正则匹配，不区分大小写

=    #进行普通字符精确匹配，与location在配置文件中的顺序无关，= 精确匹配会第一个被处理

@   #"@" 定义一个命名的 location，使用在内部定向时，例如 error_page, try_files

^~   标识符后面跟一个字符串。表示普通字符匹配，如果该选项匹配，只匹配该选项，不匹配别的选项，Nginx将在这个字符串匹配后停止进行正则表达式的匹配（location指令中正则表达式的匹配的结果优先使用），如：location ^~ /images/，你希望对/images/这个目录进行一些特别的操作，如增加expires头，防盗链等，但是你又想把除了这个目录的图片外的所有图片只进行增加expires头的操作，这个操作可能会用到另外一个location，例如：location ~* \.(gif|jpg|jpeg)$，这样，如果有请求/images/1.jpg，nginx如何决定去进行哪个location中的操作呢？结果取决于标识符^~，如果你这样写：location /images/，这样nginx会将1.jpg匹配到location ~* \.(gif|jpg|jpeg)$这个location中，这并不是你需要的结果，而增加了^~这个标识符后，它在匹配了/images/这个字符串后就停止搜索其它带正则的location。

例如：

```javascript
location  = / {
  # 只匹配"/".
  [ configuration A ] 
}
location  / {
  # 匹配任何请求，因为所有请求都是以"/"开始
  # 但是更长字符匹配或者正则表达式匹配会优先匹配
  [ configuration B ] 
}
location ^~ /images/ {
  # 匹配任何以 /images/ 开始的请求，并停止匹配 其它location
  [ configuration C ] 
}
location ~* \.(gif|jpg|jpeg)$ {
  # 匹配以 gif, jpg, or jpeg结尾的请求. 
  # 但是所有 /images/ 目录的请求将由 [Configuration C]处理.   
  [ configuration D ] 
}
```

复制

请求URI例子:

-  / -> 符合configuration A 
-  /documents/document.html -> 符合configuration B 
-  /images/1.gif -> 符合configuration C 
-  /documents/1.jpg ->符合 configuration D 

=    表示精确的查找地址，如location = /它只会匹配uri为/的请求，如果请求为/index.html，将查找另外的location，而不会匹配这个，当然可以写两个location，location = /和location /，这样/index.html将匹配到后者，如果你的站点对/的请求量较大，可以使用这个方法来加快请求的响应速度。

@    表示为一个location进行命名，即自定义一个location，这个location不能被外界所访问，只能用于Nginx产生的子请求，主要为error_page和try_files。

（2）location 优先级官方文档

1.  =前缀的指令严格匹配这个查询。如果找到，停止搜索。 
2.  所有剩下的常规字符串，最长的匹配。如果这个匹配使用^〜前缀，搜索停止。 
3.  正则表达式，在配置文件中定义的顺序。 
4.  如果第3条规则产生匹配的话，结果被使用。否则，如同从第2条规则被使用。 

（3）正则语法

~    为区分大小写的匹配。

~*   不区分大小写的匹配（匹配firefox的正则同时匹配FireFox）。

!~   不匹配的

!~*   不匹配的

.   匹配除换行符以外的任意字符

\w   匹配字母或数字或下划线或汉字

\s   匹配任意的空白符

\d   匹配数字

\b   匹配单词的开始或结束

^   匹配字符串的开始

$   匹配字符串的结束

\W   匹配任意不是字母，数字，下划线，汉字的字符

\S   匹配任意不是空白符的字符

\D   匹配任意非数字的字符

\B   匹配不是单词开头或结束的位置

捕获   (exp)   匹配exp,并捕获文本到自动命名的组里

(?<name>exp)   匹配exp,并捕获文本到名称为name的组里，也可以写成(?'name'exp)

(?:exp)   匹配exp,不捕获匹配的文本，也不给此分组分配组号

零宽断言   (?=exp)   匹配exp前面的位置

(?<=exp)   匹配exp后面的位置

(?!exp)   匹配后面跟的不是exp的位置

(?<!exp)   匹配前面不是exp的位置

注释   (?#comment)   这种类型的分组不对正则表达式的处理产生任何影响，用于提供注释让人阅读