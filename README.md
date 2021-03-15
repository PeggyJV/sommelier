# Sommelier 3.0 Stargate 

Sommelier 3.0 is a test version of the Stargate binaries.

[![codecov](https://codecov.io/gh/peggyjv/sommelier/branch/main/graph/badge.svg)](https://codecov.io/gh/peggyjv/sommelier)
[![Go Report Card](https://goreportcard.com/badge/github.com/peggyjv/sommelier)](https://goreportcard.com/report/github.com/peggyjv/sommelier)
[![license](https://img.shields.io/github/license/peggyjv/sommelier.svg)](https://github.com/peggyjv/sommelier/blob/main/LICENSE)
[![LoC](https://tokei.rs/b1/github/peggyjv/sommelier)](https://github.com/peggyjv/sommelier)
[![GolangCI](https://golangci.com/badges/github.com/peggyjv/sommelier.svg)](https://golangci.com/r/github.com/peggyjv/sommelier)


## Talk to us!

We have active, helpful communities on Twitter, Discord, and Telegram.

* [Telegram](https://t.me/peggyvaults)

## The Plan:

`sommelier` chain will consist of the `gaia` modules as well as the following custom modules:
- [`x/oracle`](https://github.com/peggyjv/sommelier/x/oracle) - A price oracle forked from the [Terra](https://terra.money) price oracle. The base denom for all prices is change from `luna` to `usd` and much of the internal code has been refactored to use more standard SDK types. This oracle will bring in price data to support our usecases (@jackzampolin)
    - [x] Import and stargate migration
    - [x] Refactors to simplify internal code for more extensibility/stability
- [`x/peggy`](https://github.com/althea-net/peggy/module/x/peggy) - An ethereum bridge that will allow assets to move to/from peggy/eth. (@jkilpatr @jackzampolin @jtremback)
    - [x] [Stargate migration](https://github.com/althea-net/peggy/pull/120)
    - [x] Import module
- [`x/il`](https://github.com/peggyjv/sommelier/x/il) - A module that consumes `oracle.GetPrice()` and `peggy.SendEthMsg()` to offer stop loss protection for LP shares (@jackzampolin @zmanian)
    - [x] Write up userflows (WIP)
    - [x] Formalize module design
    - [x] Code up module
    - [ ] Test against Rinkby
- [`frontend`](https://github.com/PeggyJV/frontend) @kkennis @jackzampolin @zmanian     
    - [x] Wireframes
    - [x] Product Requirements Doc
    - [x] Prototype
