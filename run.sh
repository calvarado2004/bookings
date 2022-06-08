#!/bin/bash

export DB_SERVER=localhost
export DB_PORT=5432
export DB_NAME=postgres
export DB_USER=postgres
export DB_PASSWORD=

go build -o bookings cmd/web/*.go && ./bookings
