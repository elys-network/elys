# Network Guide

This section provides a step-by-step guide on how to launch a new network, such as a testnet, for Elys. The guide includes instructions on how to use Ignite commands to set up and configure the new network.

<details>
<summary>Click to expand/collapse</summary>

## Coordinator Configuration

To publish the information about Elys chain as a coordinator, run the following command:

```
ignite network chain publish github.com/elys-network/elys --tag v0.1.0 --chain-id elystestnet-1 --account-balance 10000000000uelys
```

## Validator Configuration

This documentation presupposes the validator node is currently operational on `Ubuntu 22.04.2 LTS`.

### Prerequisites

Before launching a validator node, a set of tools must be installed.

To install the `build-essential` package, enter the following command:

```
sudo apt install build-essential
```

Install `go` version `1.20`

```
cd /tmp
wget https://go.dev/dl/go1.20.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
```

Append the following line to the end of the `~/.bashrc` file:

```
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
```

Run the following command:

```
go version
```

This should return the following output:

```
go version go1.20.6 linux/amd64
```

Install `ignite-cli`

Enter the following command to install the `ignite-cli` command:

```
curl https://get.ignite.com/cli! | bash
```

Then run the following command:

```
ignite network
```

Install the latest version of Elys binary by running the following command:

```
curl https://get.ignite.com/elys-network/elys@latest! | sudo bash
```

Enter the following command to initialize the validator node and request to join the network:

```
ignite network chain init 12
ignite network chain join 12 --amount 95000000uelys
```

The coordinator will then have to approve the validator requests with the following commands:

```
ignite network request list 12
ignite network request approve 12 <REQUEST_ID>,<REQUEST_ID>
```

Once all the validators needed for the validator set are approved, to launch the chain use the following command:

```
ignite network chain launch 12
```

Each validator is now ready to prepare their nodes for launch by using this command:

```
ignite network chain prepare 12
```

The output of this command will show a command that a validator would use to launch their node such as:

```
elysd start --home $HOME/spn/12 2> elysd.log &
```

A systemd service can be created to auto-start the `elysd` service.

Create the new file `/etc/systemd/system/elysd.service` with this content:

```
[Unit]
Description=Elysd Service
Wants=network.target
After=network.target

[Service]
Environment=HOME=/home/ubuntu
Type=simple
Restart=on-failure
WorkingDirectory=/home/ubuntu
SyslogIdentifier=elysd.user-daemon
ExecStart=/home/ubuntu/go/bin/elysd start --home spn/12 2>&1
ExecStop=/usr/bin/pkill elysd

[Install]
WantedBy=multi-user.target
```

Then you can use those commands to enable and start the service:

```
sudo systemctl enable elysd.service
sudo systemctl start elysd.service
```

You can check the status of the service at any time using this command:

```
sudo systemctl status elysd.service
```

Or follow the service logs by using this command:

```
sudo journalctl -u elysd.service -f
```

</details>