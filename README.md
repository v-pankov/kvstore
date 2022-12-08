# kvstore
Simple key-value store

How to run:

1. Build binaries by running `make`.

2. Setup environment by executing either `make devenv-up` or `make devenv-up-d`. The latter starts *memcached* docker container in background.

3. Execute `build/kvstore-server` in one window.

4. Execute `build/kvstore-client` in another window.