#!/bin/bash

cd sql/schema || exit 1

# Command 1
echo "Running first command down"
goose postgres postgres://postgres:postgres@localhost:5432/chirpy down-to 0

# Command 2
echo "Running second command up"
goose postgres postgres://postgres:postgres@localhost:5432/chirpy up

cd ../.. || exit 1