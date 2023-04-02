https://github.com/scylladb/scylla-code-samples/tree/master/mms

### コンテナの起動
```shell
docker-compose build

docker-compose up -d
```

### コンテナの後片付け
```shell
docker-compose down 
```

### 生存確認
```shell
docker exec -it cassandra_node1 nodetool status
docker exec -it cassandra_node2 nodetool status
```