# GoFish

[![Build Status](https://travis-ci.org/thenomemac/gofish.svg?branch=master)](https://travis-ci.org/thenomemac/gofish)

Golang refactor of https://github.com/thenomemac/fishapp


## launching app on server

``` bash
docker run -d -v ~/.aws:/root/.aws -p 80:3000 thenomemac/gofish
```
