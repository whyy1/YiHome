docker run  --name mysql  -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
docker cp  mysql:/var/lib/mysql/ e:\volumes\mysql\
docker run -dit --name mysql  -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7 -v /e/volumes/mysql/:/var/lib/mysql/
docker run -dit --name mysql  -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7  -v /e/volumes/mysql/mysql-files:/var/lib/mysql-files -v /e/volumes/mysql/log:/var/log/mysql -v /e/volumes/mysql/data:/var/lib/mysql -v /e/volumes/mysql/conf:/etc/mysql
docker run -p 3306:3306 --name mysql -v /f/docker/mysql8/mysql-files:/var/lib/mysql-files -v /e/volumes/mysql/log:/var/log/mysql  -v /e/volumes/mysql/data:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=root mysql:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
//consul启动
docker run -d --network=host --name=consul consul:latest agent -server -bootstrap -ui -node=1 -client='0.0.0.0' -bind='172.29.251.146'
//微服务user启动
docker run --net=host -v ~/conf/:/conf --name user registry.cn-heyuan.aliyuncs.com/whyy1/yihome:user-v1
//微服务getCaptcha启动
docker run --net=host -v ~/conf/:/conf --name getCaptcha registry.cn-heyuan.aliyuncs.com/whyy1/yihome:getCaptcha-v1
//redis启动
sudo docker run -p 63799:6379   --name redis -v ~/volumes/redis/conf:/etc/redis/redis.conf -v ~/volumes/redis/data:/data -d redis redis-server /etc/redis/redis.conf --appendonly yes --requirepass "heikebisi"
//进入redis客户端
docker exec -ti redis redis-cli -h 127.0.0.1  -p 6379
//mysql启动
docker run -p 33069:3306 --name mysql --restart=always -d -v ~/volumes/mysql/conf/my.cnf:/etc/mysql/my.cnf -v ~/volumes/mysql/mysql-files:/var/lib/mysql-files -v ~/volumes/mysql/log:/logs -v ~/volumes/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=kTspuYb8CCI3aplU mysql:5.7  --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci