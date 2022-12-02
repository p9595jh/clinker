# ganache

Clinker is using [ganache](https://github.com/trufflesuite/ganache) to test.

To start, you have to make `.env` file and fill `PK1`, `PK2` with your private keys.

This repository provides methods to use `ganache`, which are `cli` and `docker`.

`cli` needs `ganache` command which can be installed from below:

```
npm install ganache --global
```

Then you can use CLI commands.

```
# to start
./cli.sh start

# to stop
./cli.sh stop
```

`docker` method can be run by docker-compose, so docker and docker-compose must be installed first.

```
docker-compose -f docker-compose.yml up -d
```
