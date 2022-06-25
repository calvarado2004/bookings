#!/bin/sh

#This script will be running on the init container

export ADMIN_DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@$DB_SERVER:$DB_PORT/postgres
export DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@$DB_SERVER:$DB_PORT/$DB_NAME 

sed -i "s/{{ DB_USER }}/$DB_USER/g" admin_migrations/20220509161150_database.postgres.up.sql
sed -i "s/{{ DB_PASSWORD }}/$DB_PASSWORD/g" admin_migrations/20220509161150_database.postgres.up.sql

sed -i "s/{{ DB_USER }}/$DB_USER/g" admin_migrations/20220509161150_database.postgres.down.sql
sed -i "s/{{ DB_PASSWORD }}/$DB_PASSWORD/g" admin_migrations/20220509161150_database.postgres.down.sql


/usr/local/bin/soda migrate -e admin -p ./admin_migrations status 
/usr/local/bin/soda migrate -e admin -p ./admin_migrations up 
/usr/local/bin/soda migrate -e admin -p ./admin_migrations status 

/usr/local/bin/soda migrate -e production status
/usr/local/bin/soda migrate -e production up
/usr/local/bin/soda migrate -e production status
