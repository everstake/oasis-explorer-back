# Oasis Explorer

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Oasis Explorer is an open-source Tezos explorer based on own indexer core.
Developed and supported by Everstake team.

The indexer fetches raw data from the Oasis node, then processes it and stores in the database in such a way as to provide effective access to the blockchain data. 
None of this can be accessed via node RPC, but Oasis Explorer indexer makes this data (and much more) available.

Oasis Explorer provides a REST-like API, so you don't have to connect to the database directly.
Our API server will do for you all needed aggregations.

Full documentation you can find here [OasisExplorerAPI](https://oasismonitor.com/docs) 
## Getting started

### Dependencies

#### Clickhouse
 To install Clickhouse follow guide below
 
https://clickhouse.tech/docs/en/getting-started/install/

 Than create an empty database and its user
#### Postgres
 To install Postgres follow guide below  
 
 https://www.postgresqltutorial.com/postgresql-getting-started/
 
 Than create an empty database and its user
#### Oasis-node
 To run a Non-validator Node follow steps below
 
 https://docs.oasis.dev/general/run-a-node/set-up-your-node/run-non-validator
### Installing and running oasis-wallet

```
git clone https://github.com/everstake/oasis-explorer.git
cd oasis-explorer/
mkdir .secrets
```

Download latest `genesis` file

Mainnet
```
wget https://github.com/oasisprotocol/mainnet-artifacts/releases/download/2021-04-28/genesis.json
``` 

Testnet

```
wget https://github.com/oasisprotocol/testnet-artifacts/releases/download/2021-04-13/genesis.json
```

Setup `config.json` into `.secrets` folder.

Run
```
docker-compose up --build -d 
```

Then curl health status of API
```
curl http://localhost:9000/health
```

### OpenAPI
Oasis-Explorer exposes an [OpenAPI](https://github.com/everstake/oasis-explorer/blob/master/swagger/swagger.yml)

Could be found at http://localhost:8080
