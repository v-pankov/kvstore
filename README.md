# kvstore
Simple key-value store

How to run:

1. Generate grpc schema by running `make generate`.

2. Run development environment by executing either `make devenv-up` or `make devenv-up-d`. The latter starts *memcached* docker container in background.

3. Open another terminal and run `make run-memcached` to execute *memcached* library client example program.

4. Open another ternimal and run `make run-server` to start key-value store server.

5. Open another terminal and run `make run-client` to execute key-value store client example program.