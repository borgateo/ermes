```
███████╗██████╗ ███╗   ███╗███████╗███████╗
██╔════╝██╔══██╗████╗ ████║██╔════╝██╔════╝
█████╗  ██████╔╝██╔████╔██║█████╗  ███████╗
██╔══╝  ██╔══██╗██║╚██╔╝██║██╔══╝  ╚════██║
███████╗██║  ██║██║ ╚═╝ ██║███████╗███████║
╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝
```
## Features

1. Follow all(*) the followers of a passed user
2. Like the latest 3 (configurable) posts of all your followers
3. Like the latest 3 (configurable) posts of all your followings
4. Like the latest ~16 posts of your timeline
5. Unfollow people that don't follow you back
6. More coming soon...

(*) Note: only the accounts that meet a specific criteria.

## Install
`$ go get -v github.com/borteo/ermes`

## Get started
If you have already govendor installed, start from point 2:

1. `$go get -u github.com/kardianos/govendor`
2. `$ govendor sync`
3. `$ go build`
4. `$ ./ermes [-h -followers -followings -user={username} -timeline] [-skip]`

Use the flag `-h` to see all the available flags.

For instance, to run feature #5:
`$ ./ermes -unfollow`

Note: `-skip` flag works only with followers, followings and user


## Contributing
This is a WIP project, feel free to fork it and/or create PRs/Issues.

## Shout-out
- https://github.com/ahmdrz/goinsta
- https://github.com/kirsle/follow-sync


## Disclaimer

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Instagram or any of its affiliates or subsidiaries. This is an independent and unofficial pseudo-bot. Use at your own risk.
