# Black Rocket, monitoring tools for Cardano blockchain nodes

Black Rocket Informer 

## How to support this project

If you are Cardano enthusiast and/or SPO, you can support us by delegating your ADA to our node

https://pooltool.io/pool/59b59f232e80f18ce64e0f74560effbf49a3e95ddf6d079681db8b23/epochs

Or you can do it directly by sending coins to us, our wallets:

#### Cardano: 
addr1qyqa4x78s0l9vusy3kw4772najwzer2s0pk9l8t4hfrushsm50kqnkgssjre0nysnwz9uc20gsanmqsnwxdnxj4w7zfswl9fse

#### Ethereum or Binance Smart chain: 
0x1af0637A6131f29389c2e68517D61bF5e2655a57

If you are participating in Project Catalyst, you can vote for this solution or comment it on the ideascale: 
https://cardano.ideascale.com/a/dtd/Monitoring-solution-for-a-node/383557-48088

## How to install Informer

### Installing directly from the binaries

Upcoming, not available yet due to active development.

Binaries will be abailable for Linux only.

### Installing Go
Controller requires Go 1.13 to compile, please refer to the [official documentation](https://go.dev/doc/install) for how to install Go in your system.

### Installing Informer:
```
go get github.com/adarocket/informer 
```

### Сonfiguration 
The next step is to create a configuration file.
The program may promt you to create a configuration file or you can create it manually.
* ticker - name of the informer.
* uuid - unique id of the informer.
* location - 
* controller_url - 

```
cat << EOF > ~/etc/ada-rocket/informer.conf
{
    "ticker"                        : "[BKV]",
    "uuid"                          : "927b9103-aa64-44be-99c6-85e55695ea50",
    "location"                      : "Belarus",
    "node_monitoring_url"           : "http://127.0.0.1:12788",
    "controller_url"                : "165.22.92.139:5300",
    "time_for_frequently_update"    : 10,
    "time_for_rare_update"          : 60,
    "blockchain"                    : "cardano"
}
EOF
```


### Black Rocket Community Support

Feel free to join our community channel on telegram: https://t.me/+R6gKUjtqMVY1OTBi



