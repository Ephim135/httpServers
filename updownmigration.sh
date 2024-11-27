#!/bin/bash

cd sql/schema || exit 1

# Command 1
echo "Running first command"
goose postgres postgres://postgres:postgres@localhost:5432/chirpy down

# Command 2
echo "Running second command"
goose postgres postgres://postgres:postgres@localhost:5432/chirpy up

cd ../.. || exit 1