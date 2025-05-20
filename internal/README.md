## Backend README

## Running the app

Considering you are in the root of the project directory:

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
go run {proj_root}/internal/main.go &
```

In order for the server to run, you must have valid openssl certificates in *{proj_root}/certs*. To create them, use this command in *proj_root*:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/server.key -out certs/server.crt -config certs/openssl.cnf -extensions v3_req
```

### Tests

```bash
cd {proj_root}/internal/test/db_test && go test -v; cd -
cd {proj_root}/internal/test/server_test && go test -v; cd -
```
