# lolachain
A blockchain named for my dog, Lola. Currently two tokens exist on the chain, RockyCoin (RKY) and LolaCoin (LOLA). This is largely a personal project done for fun.

Currently a work in progress, I plan to implement a decentralized consensus model among some other features.

## Building
A Makefile is provided at the root of the project, for convenience. Targets exist for Windows, Linux, and macOS. Simply run `make linux`, `make darwin`, or `make windows`. Targets to `test` and `lint` also exist. Files are written to platform specific subdirectories in `./dist/`.

## Running
lolachain consists of three applicaions: `lolachain-validator`, `lolachain-gui`, and `lolachain-wallet`. The main application for mining is `lolachain-validator`, and can simply be run via `./dist/linux/lolachain-validator`. The GUI is a bit of a work in progress, and can be run via `./dist/linux/lolachain-gui`. Then, open a browser to `http://localhost:8080`. The `lolachain-wallet` application is simply a CLI wallet. At the time of first launch, the applications will create a local wallet for you via a private key file at `~/.lolachain/key.pem`.

## TODO
- [ ] Test Coverage, there's some but not nearly enough
- [ ] Rewrite UI to not be client/server based
- [ ] Add decentralized concensus model (lol kind of a big hole)
- [ ] Refactor CLI wallet