#!/bin/sh

echo "Waiting for postgres..."
sleep 2

while ! nc -z $1 $2; do
  sleep 0.1
done

echo "Postgres SQL started"
