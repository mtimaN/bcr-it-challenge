CREATE ROLE IF NOT EXISTS admin WITH PASSWORD = 'CPass0319' AND SUPERUSER = true AND LOGIN = true;

CREATE ROLE IF NOT EXISTS backend WITH PASSWORD = 'BPass0319' AND LOGIN = true;

CREATE KEYSPACE IF NOT EXISTS cass_keyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

GRANT ALL PERMISSIONS ON KEYSPACE cass_keyspace TO backend;


CREATE TABLE IF NOT EXISTS cass_keyspace.users (
	username text PRIMARY KEY,
	password text,
	email text,

	-- data
	category int
);
