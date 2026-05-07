#!/bin/bash
# 初始化 MySQL 数据库

echo "Waiting for MySQL to be ready..."
sleep 5

# 创建数据库和表
mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < /docker-entrypoint-initdb.d/yuanzi.sql

echo "Database initialized successfully!"
