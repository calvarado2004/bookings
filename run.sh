#!/bin/bash

export DB_SERVER=localhost
export DB_PORT=5432
export DB_NAME=bookings
export DB_USER=postgres
export DB_PASSWORD=postgres
export MAILHOG_HOST=localhost
export MAILHOG_PORT=1025

export DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@$DB_SERVER:$DB_PORT/$DB_NAME

go build -o bookings cmd/web/*.go && ./bookings
