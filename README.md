# MDB

> A CLI for better management of MongoDB commands

## Overview

I got tired of copy-pasting urls from different environments while running MongoDB commands.

## Commands

This CLI supports the following MongoDB commands:

- `mongoimport`: import data to a MongoDB collection
- `mongoexport`: export data from a MongoDB collection
- `mongodump`: create a binary export from MongoDB
- `mongorestore`: import data from a MongoDB dump

Aditional commands:

- `setenv`: to set an environment for each db connection (dev, staging, prod, etc.)

## How to run

Make sure you have [Docker](https://www.docker.com/products/docker-desktop) installed in your machine. The commands are executed inside a Docker container, so the CLI doesn't depend on `mongo` shell local installations.

## Pending

- [ ] improve code documentation
- [ ] add [go-prompt](https://github.com/c-bata/go-prompt) to auto-complete commands
