#!/usr/bin/env bash
set -e

source local.env

# specify function that creates database with given name

psql << EOF
DROP DATABASE IF EXISTS ${SVC_DB_NAME}_test;
CREATE DATABASE ${SVC_DB_NAME}_test;

EOF

psql << EOF
DROP DATABASE IF EXISTS ${SVC_DB_NAME};
CREATE DATABASE ${SVC_DB_NAME};

EOF

goose -dir migration up
export GOOSE_DBSTRING="user=postgres sslmode=disable port=5432 password=postgres host=localhost dbname=unit_test"
goose -dir migration up


