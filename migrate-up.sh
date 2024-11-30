#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | xargs)
else
    echo ".env file not found!"
    exit 1
fi

# Run goose migration (up)
goose -dir ./sql/schema postgres $DATABASE_URL up

