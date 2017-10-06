
## Install
`$ go get -v github.com/borteo/ermes`

## Run w/o docker-compose

If you have already govendor install, start from point 2:

1. `$go get -u github.com/kardianos/govendor`
2. `$ govendor sync`
3. `$ go build`
4. `$ ./ermes`


## Run w/ docker-compose
### Requirements
Docker
Docker-compose

### How to run
`docker-compose up --build`

Refs:

- https://github.com/ahmdrz/goinsta
- https://github.com/kirsle/follow-sync


