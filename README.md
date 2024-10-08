# debezium latency test

```
Mysql > Debezium > Kafka
```

## Usage

```
git clone https://github.com/alok87/dbz
cd dbz
docker-compose up --build
```

- producer will start posting msgs every 10 seconds
- consumer will show the time it received

Kafka, zookeeper, debezium, mysql, producer and consumer starts using docker-compose.
