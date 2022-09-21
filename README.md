# 使用说明

## 1 执行方法

---
install.sh是基础组件集成安装脚本

执行方法: `bash install.sh <mode> <module>`

其中mode的可选项有：

- single：组件单机安装
- cluster： 组件集群安装

其中module可选项有：

- mysql
- redis
- mongodb
- zookeeper
- rabbitmq
- apollo
- lts
- memcache
- nginx

组件安装前先执行download_package.sh下载组件安装包
> bash download_package.sh

用户密码随机生成保存在`install.config`和`base_component.properties`文件，重复执行密码会覆盖

## 2 组件说明

---

### mysql

    1. mysql支持单机安装和集群安装，集群模式是主从集群
    2. 单机安装只需修改mysql_host的ip地址
    3. 集群安装需修改mysql_num和mysql_cluster配置，mysql_cluster{num}数量和mysql_num的值相同
    4. 支持一主多从，默认mysql_cluster的第一个ip为主
    5. 默认安装root、admin、app、repl四个用户，用户密码安装时自动生成

### redis

    1. redis支持单机安装和哨兵集群安装
    2. 单机安装只需修改redis_host的ip地址
    3. 集群安装填写redis_num和redis_cluster配置，redis_cluster{num}数量和redis_num相同
    4. 支持一主多从从，默认redis_cluster填写的第一个ip为主
    5. redis_cluster_name不能修改，固定为redisprd

### mongodb

    1. mongdodb支持单机安装和集群安装，集群为分片模式
    2. 单机安装只需修改mongdodb_host的ip地址
    3. 集群安装填写mongdodb_num和mongdodb_cluster，mongdodb_cluster{num}数量和mongdodb_num相同
    4. 默认安装root、app三个用户，用户密码安装时自动生成

### elasticsearch

    1. elasticsearch支持单机安装和集群安装
    2. 单机安装需要修改elasticsearch_host配置
    3. 集群安装填写elasticsearch_num和elasticsearch_cluster{num}配置，集群节点数不限
    4. elasticsearch内存设置为服务器内存一半

### rabbitmq

    1. rabbitmq支持单机安装和集群安装，集群为镜像模式
    2. 单机安装需要修改rabbitmq_host配置
    3. 集群安装填写rabbitmq_num和rabbitmq_cluster{num}配置，集群节点数不限
    4. 脚本默认已经安装delayed_message延时插件

### zookeeper

    1. zookeeper支持单机安装和集群安装
    2. 单机安装需要修改zookeeper_host配置
    3. 集群安装填写zookeeper_num和zookeeper_cluster{num}配置，集群节点数不限

### memcache

    1. memcache支持单机安装
    2. 单机安装需要修改memcache_host配置

### apollo

    1. apollo支持单机安装、多环境集群安装
    2. 单机安装支持pro和fat环境安装，设置apollo_env=pro或者apollo_env=fat
    3. 集群安装可以设置apollo_env=pro,uat或者apollo_env=pro，apollo_env=pro,uat需要3台服务器，apollo_env=pro需要2台服务器
    4. 数据库根据环境不同分为生产数据库、A服务数据库、测试数据库
        apollo_root_pwd：数据库root账号密码
        apollo_prod_pwd：生产数据库账号密码，账号为apollo_prod_user
        apollo_a_pwd： A服数据库账号密码，账号为apollo_a_user
        apollo_test_pwd: 测试数据库账号密码，账号为apollo_test_user
    5. 数据库默认安装在apollo环境第一台服务器

### lts(Light Task Schedule)

    1. lts支持单机安装
    2. lts数据库默认与apollo共用一台数据库，安装lts前需先安装apollo
        lts_root_pwd：数据库root账号密码，填写apollo数据库root账号密码
        lts_jdbc_url：数据库地址，格式为："jdbc:mysql://192.168.253.101/lts_prod_db"
        lts_jdbc_ip： 数据库ip
        lts_jdbc_user: lts数据库用户
        lts_jdbc_pwd：lts数据库密码，自动生成
    3. lts_env可以设置lts环境为prod或者test环境
        lts_env=prod时，需设置lts_jdbc_url中DB名称为lts_prod_db
        lts_env=test时，需设置lts_jdbc_url中DB名称为lts_test_db
    4. lts_zookeeper_url设置zookeeper地址，多个节点以逗号(,)分隔
        格式：zookeeper://192.168.253.101:2181,192.168.253.102:2181,192.168.253.103:2181

### nginx

    1. nginx只支持代理后端服务
    2. 修改nginx配置请修改nginx.conf文件
    3. 在upstream.txt文件中填写后端服务ip、端口
    4. 如需增加upstream，需在nginx.conf和upstream.txt文件中添加对应的信息
