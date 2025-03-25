#!/usr/bin/env bash

# Install migrate CLI
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

echo "✅ migrate CLI installed"
