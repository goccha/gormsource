# GORM用DataSourceライブラリ

# examplesの実行
## RDBの起動
- シングル構成
```shell
> cd deployments
> docker-compose up -d
```

- primary-replica構成
```shell
> cd deployments
> docker-compose -f docker-compose-replication.yml up --detach --scale mysql-master=1 --scale mysql-slave=1 --scale postgresql-master=1 --scale postgresql-slave=1
```

