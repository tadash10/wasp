![Wasp logo](documentation/static/img/logo/WASP_logo_dark.png)

#  Legacy-Migration Wasp

For the regular wasp node software go to https://github.com/iotaledger/wasp .

<!-- TODO add link -->
This version of the wasp software was made with the objective of allowing the [migration of legacy funds](.) from the pre-chysalis network.

Contains a simple [contract](./packages/legacymigration/interface.go) that can be called to release funds given a valid signature of an unmigrated bundle.

This repo also contains a [snapshot](./packages/legacymigration/migratable.csv) of the old network containing all the unmigrated bundles

Everytime a migration is successful, the funds are released on L1 to the target address and an [event](./packages/legacymigration/impl.go:102) is published, making the entire process auditable.

<!-- TODO add link -->
At any point, the [governance contract of the EVM Chain](.) can vote and decide to burn the unmigrated tokens.

## Instructions

<!-- TODO -->
All committee participants must use the [node-docker-setup with the legacy-migration wasp software](.).

The deployer must use the [wasp-cli compiled from this branch]().

      - deploy command
      - funds must be deposited to the chain
      - call views to confirm everything is as expected 
      - relinquish gov controller of the chain UTXO
