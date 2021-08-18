# Oasis Explorer

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Oasis Explorer is an open-source Oasis explorer based on its own indexer core.
Developed and supported by the Everstake team.

The indexer fetches raw data from the Oasis node, then processes and
stores it in the database in such a way 
as to provide effective access to the blockchain data.
None of this can be accessed via node RPC,
but Oasis Explorer indexer makes this data (and much more) available.

Oasis Explorer provides a REST-like API, 
so you don't have to connect to the database directly.
Our API server will do all the needed aggregations for you.

Full documentation you can found here [OasisExplorerAPI](https://oasismonitor.com/docs)

## Getting started

### Dependencies

#### Clickhouse

 To install Clickhouse follow the guide below

<https://clickhouse.tech/docs/en/getting-started/install/>

 Then create an empty database and its users

#### Postgres

 To install Postgres follow the guide below  

 <https://www.postgresqltutorial.com/postgresql-getting-started/>

 Then create an empty database and its users

#### Oasis-node

 To run a Non-validator Node follow the steps below

 <https://docs.oasis.dev/general/run-a-node/set-up-your-node/run-non-validator>

### Installing and running oasis-wallet

```bash
git clone https://github.com/everstake/oasis-explorer.git
cd oasis-explorer/
mkdir .secrets
```

Download latest `genesis` file

Mainnet

```bash
wget https://github.com/oasisprotocol/mainnet-artifacts/releases/download/2021-04-28/genesis.json
```

Testnet

```bash
wget https://github.com/oasisprotocol/testnet-artifacts/releases/download/2021-04-13/genesis.json
```

Setup `config.json` into `.secrets` folder.

Run

```bash
docker-compose up --build -d 
```

Then curl health status of API

```bash
curl http://localhost:9000/health
```

### OpenAPI

Oasis-Explorer exposes an [OpenAPI](https://github.com/everstake/oasis-explorer/blob/master/swagger/swagger.yml)

Could be found at <http://localhost:8080>
