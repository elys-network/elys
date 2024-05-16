# Elys

**Elys** is a blockchain built using Cosmos SDK and CometBFT. It is designed to be a fast, scalable, and secure blockchain that can be used to build decentralized applications.

## Installation

### With Makefile (Recommended)

This section provides a step-by-step guide on how to build the Elys Chain binary from the source code using the provided makefile. The makefile automates the build process and generates a binary executable that can be run on your local machine.

1. Clone the Elys chain repository:

```bash
git clone https://github.com/monopauli/elys.git
```

2. Navigate to the cloned repository:

```bash
cd elys
```

3. Optionally, checkout the specific branch or tag you want to build:

```bash
git checkout <version>
```

4. Ensure that you have the necessary dependencies installed. For instance, on Ubuntu you need to install the `make` tool:

```bash
sudo apt-get install --yes make
```

5. Run the `make build` command to build the binary:

```bash
make build
```

6. The binary will be generated in the `./build` directory. You can run the binary using the following command:

```bash
./build/elysd
```

You can also use the `make install` command to install the binary in the `bin` directory of your `GOPATH`.

## Release

To release a new version of Elys, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

## Learn more

- [Twitter](https://twitter.com/elys_network)
- [TestNet Explorer](https://testnet.ping.pub/elys)
- [Developer Chat](https://discord.gg/3JtgtGJ3By)
- [Github](https://github.com/elys-network)
