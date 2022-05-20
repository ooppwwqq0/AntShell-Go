# AntShell 说明

```
                 _______  __    _  _______  _______  __   __  _______  ___      ___
                |   _   ||  |  | ||       ||       ||  | |  ||       ||   |    |   |
                |  |_|  ||   |_| ||_     _||  _____||  |_|  ||    ___||   |    |   |
                |       ||       |  |   |  | |_____ |       ||   |___ |   |    |   |
                |       ||  _    |  |   |  |_____  ||       ||    ___||   |___ |   |___
                |   _   || | |   |  |   |   _____| ||   _   ||   |___ |       ||       |
                |__| |__||_|  |__|  |___|  |_______||__| |__||_______||_______||_______|


                               ___           ___              _   _      _
                              / __\_   _    / __\__ _ ___ ___| |_(_) ___| |
                             /__\// | | |  / /  / _` / __/ __| __| |/ _ \ |
                            / \/  \ |_| | / /__| (_| \__ \__ \ |_| |  __/ |
                            \_____/\__, | \____/\__,_|___/___/\__|_|\___|_|
                                   |___/
```

## 简介

> 声明：本人因能力问题，代码质量差，读者见谅

`AntShell`（文本终端连接工具）不与其他类似工具（如：`ansible`）对标，主要用于个人工作、学习过程中快速、高效管理服务器维护工作，之所以选择自己开发而没有使用开源软件有以下几点原因：
1.命令行工具结合一些过滤条件更加高效
2.数据存储完全存储在本地，数据不会被某些无良软件上传导致泄露
3.实战项目的开发过程中能够更好的学习编程知识
4.自己的习惯和需求自己解决，自己就是用户更能感受用户的实际需求

### 功能简介

#### 已有功能

* 添加、删除、修改主机记录
* 多种搜索记录策略
* 主机自动登录
* 堡垒机模式登录主机（目前已测试齐治堡垒机）
* 登录过程中设置中文环境变量，自动切换用户、执行命令
* 登录自动工作目录配置
* 常用主机优先级
* 本地主机记录排序
  * 根据排序字段进行排序展示
* 多个banner可选，有default策略、random策略
* 堡垒机密码失效等待（堡垒机秘钥会有一段时间的有效期，在有效期内只能使用一次）

#### 未来功能列表

* 本地数据备份
* 获取主机信息插件
* 新的UI界面：webUI、TUI，GUI
* 批量命令、文件上传、文件下载
* 自定义banner

### `AntShell`执行命令的由来

执行命令为`a`：是因为字母`a`位于键盘中手指最容易按到的位置，为了让使用者能够更便捷、更快速的使用，因此只使用了一个字母作为命令输入。

### `AntShell`发展历史

> 2016年由shell语言编写执行命令命名为a，同年起名AutoSSH
> 
> 2017年初使用python语言重新编写，改名为Adam，取自希腊神话人物名：亚当
> 
> 2017年9月14日正式更名AntShell
> 
> 2021年8月9日使用Golang语言重写，版本号升级为1.0

## 功能参数解析

### 功能模式

#### 交互模式

* 当已知入参无法选取唯一主机记录时，进入交互模式
* 交互模式下支持：主机id选择，主机名模糊匹配，主机ip模糊匹配，主机ip精确匹配；
* 多次模糊匹配输入支持进一步筛选，支持清楚筛选记录
* 交互模式下会自动进入单列模式，支持翻页

#### 快速模式|无变量参数

* 支持一位无变量参数，无变量参数可以直接进行精确搜索、模糊搜索、指定id选定主机记录
* 通过无变量参数可以快速筛选主机记录，跳过交互模式，无变量参数可以和有变量参数配合使用

```shell
a 10.0.0.1
a app -n 1
```
### 常用参数组合

```shell
# 筛选主机
a -s name | a -s ip
a -n num
a -s name -n num
a num | a name | a ip
```

### 参数优先级

待续

## 性能提升

```bash
# python
a -l  0.32s user 0.10s system 96% cpu 0.432 total
# go
a -l  0.10s user 0.03s system 79% cpu 0.162 total
```

## HELP

```
AntShell version: AntShell/1.2
GitHub: https://github.com/ooppwwqq0/AntShell-Go
Usage: antshell|a [ -h | -version ] [-l [-m 2] ] [ v | -n 1 | -s 'ip|name' ] [ -A ] [ -B ]
        [ -e | -d ip | -sort | -a ip [-name tag | -user root | -passwd *** | -port 22 | -sudo root | -path /tmp] ]
  -B    堡垒机模式
  -a string
        添加主机信息并登陆
  -d    删除主机信息并退出
  -e    编辑主机信息
  -l    输出主机列表并退出
  -m int
        列表显示列数1-5
  -n int
        选择连接的主机编号
  -name string
        本地主机别名
  -passwd string
        密码
  -path string
        初始目录
  -port int
        端口
  -s string
        模糊匹配主机信息
  -sort
        排序
  -sudo string
        指定sudo用户
  -totp
        打印 totp
  -user string
        登录主机用户名
  -version
        打印版本信息并退出

# Add host record       
        a -a 10.0.0.1
        a -a 10.0.0.1 -name app01
        a -a 10.0.0.1 -name app01 -passwd 123456
        a -a 10.0.0.1 -name app01 -user root -passwd 123456 
        a -a 10.0.0.1 -name app01 -user root -passwd 123456 -port 22022
        a -a 10.0.0.1 -name app01 -user ubuntu -passwd 123456 -sudo root
        a -a 10.0.0.1 -name app01 -user ubuntu -passwd 123456 -path /tmp
        a -a 10.0.0.1 -name app01 -user ubuntu -passwd 123456 -port 22022 -sudo root -B
# Delete host record
        a -d 10.0.0.1
        a -d app01
# Edit host record
        a -e
        a -e -s 10.0.0.1
        a -e -s app01 -n 2
# List host record
        a -l
        a -l -m 2
# Sort host record +10
        a -sort
        a -sort -n 1
        a -sort app01 -n 1
# Login host
        a
        a 2
        a app01
        a 10.0.0.0.1
        a app01 -n 2
        a -s 10.0.0.1 -n 1
        a -s app01 -n 2
```
