# tendermint-bugs
A collection of suspected bugs found in tendermint.

Should be run with forked tendermint which is integrated to run with netrix. Can be found [here](https://github.com/zeu5/tendermint/tree/pct-instrumentation)

## Getting started
### Requirements
- Golang 1.18
- Docker
- [docker-compose](https://pypi.org/project/docker-compose/). Not the one built-in to recent versions of docker. I found pip to be the easiest way to install it: `pip3 install --user docker-compose`.

### Instrumented nodes
The tests require a modified version of tendermint found [here](https://github.com/zeu5/tendermint/tree/pct-instrumentation). 
You may use the script `third_party.sh` to download and configure it.
Alternatively, a more comprehensive guide is available [here](https://github.com/netrixframework/tendermint-testing/blob/master/README.md).

Next you need to build the docker image and create the bridge network:
```shell
cd third_party/tendermint-pct-instrumentation/

# Build the linux binary in ./build
make build-linux

# Build tendermint/localnode image
make build-docker-localnode

# Create the bridge network
docker-compose up --no-start
```

## Running the tests
Start the testing server **before** the tendermint nodes:

```shell
go run ./server.go -bug dummy
```

In another terminal, start the tendermint nodes:

```shell
cd ../tenderint-pct-instrumentation
make localnet-start
```

Eventually you should see this line in the test output:
```json
{"level":"info","msg":"Testcase succeeded","service":"TestingServer","testcase":"RoundSkipWithPrevotes","time":"2022-05-20T11:11:06+02:00"}
```

You can then stop both the server and the nodes (with Ctrl-C).

To check buggy testcases, change the `-bug` flag to something like `-bug bug001`.