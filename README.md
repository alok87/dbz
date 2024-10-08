# debezium latency test

```
Producer > Mysql > Debezium > Kafka < Consumer
```

## Usage

- Start services
```
git clone https://github.com/alok87/dbz
cd dbz
export DEBEZIUM_VERSION=2.1
docker-compose up --build
```

- Start connector
```
curl -i -X POST -H "Accept:application/json" -H  "Content-Type:application/json" http://localhost:8083/connectors/ -d @register-mysql.json
```

- producer will start posting msgs every 10 seconds
- consumer will show the time it received (mesaure the time)

Kafka, zookeeper, debezium, mysql, producer and consumer starts using docker-compose.

## Appendix

- Logs of producer
```
docker-compose logs -f producer
```

- Logs of consumer
```
docker-compose logs -f consumer
```

- May need to restart producer and consumer
```
docker-compose restart producer
docker-compose restart consumer
```

- Login to SQL
```
docker-compose -f docker-compose-mysql.yaml exec mysql bash -c 'mysql -u $MYSQL_USER -p$MYSQL_PASSWORD inventory'
```
