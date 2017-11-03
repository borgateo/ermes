
███████╗██████╗ ███╗   ███╗███████╗███████╗
██╔════╝██╔══██╗████╗ ████║██╔════╝██╔════╝
█████╗  ██████╔╝██╔████╔██║█████╗  ███████╗
██╔══╝  ██╔══██╗██║╚██╔╝██║██╔══╝  ╚════██║
███████╗██║  ██║██║ ╚═╝ ██║███████╗███████║
╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝

## Install
`$ go get -v github.com/borteo/ermes`

## Run w/o docker-compose

If you have already govendor install, start from point 2:

1. `$go get -u github.com/kardianos/govendor`
2. `$ govendor sync`
3. `$ go build`
4. `$ ./ermes`

There are several flags available when you launch ermes:

Unfollow people that don't follow you back:
`./ermes -unfollow`

Like your followers or followings media:
`	./ermes -followers`
`	./ermes -followings`

Follow passed user's followers and like their media:
`./ermes -user=username -reset=true`

## Run w/ docker-compose
### Requirements
Docker
Docker-compose

### How to run
`docker-compose up --build`

Refs:

- https://github.com/ahmdrz/goinsta
- https://github.com/kirsle/follow-sync


### Disclaimer

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Instagram or any of its affiliates or subsidiaries. This is an independent and unofficial pseudo-bot. Use at your own risk.
