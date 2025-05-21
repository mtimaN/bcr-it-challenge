## Backend README

## Running the app

To configure the project root, use

```bash
EXPORT proj_root='/path/to/project/root'
```

### Databases

```bash
docker-compose up -d
```

Wait until gossip settles...

```bash
docker exec -it cassandraDB cqlsh
```

from schema.cql:
```cqlsh
CREATE KEYSPACE IF NOT EXISTS cass_keyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS cass_keyspace.users (
	username text PRIMARY KEY,
	password text,
	email text,

	-- data
	category int
);
```

### Server

To start in background:

```bash
cd $proj_root/internal/ && go run .
```

In order for the server to run, you must have valid openssl certificates in *$proj_root/certs*. To create them, use this command:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $proj_root/certs/server.key -out $proj_root/certs/server.crt -config $proj_root/certs/openssl.cnf -extensions v3_req
```

You can press *Enter* until the setup finished.

### Tests

```bash
cd $proj_root/internal/test && go test -v; cd -
```

In case a Cassandra test fails, you might have to run:

```cqlsh
TRUNCATE cass_keyspace.users;
```

inside the docker cqlsh (see *Databases* section) before running that test manually


If the issue persists, run the tests manually :(.

## Closing the app

Close the server with a simple interrupt (*CTRL+C*), then run:

```bash
docker-compose down
```

Erasing the data:

```bash
docker volume ls
docker image ls
```

will yield a list of volumes, some of them prefixed *'bcr-it-challenge'*. Remove them with:

```bash
docker volume rm $volume_name
docker image rm $image_name
```
