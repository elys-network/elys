<!--
order: 2
-->

# Usage

## Commands

### Querying Asset Entries

```bash
ASSET_BASE_DENOM="udenom"

elysd query assetprofile show-entry $ASSET_BASE_DENOM
elysd query assetprofile list-entry
```

### Querying Module Parameters

```bash
elysd query assetprofile params
```

### Add Entry using Gov Proposal

A proposal can be submitted to add one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "add new entries",
  "description": "add new entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgCreateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2",
      "decimals": "18",
      "denom": "mytoken",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken",
      "displaySymbol": "mytoken",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken",
      "ibcCounterpartyDenom": "mytoken",
      "ibcCounterpartyChainId": "test"
    },
    {
      "@type": "/elys.assetprofile.MsgCreateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken3",
      "decimals": "18",
      "denom": "mytoken",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken",
      "displaySymbol": "mytoken",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken",
      "ibcCounterpartyDenom": "mytoken",
      "ibcCounterpartyChainId": "test"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Update Entry using Gov Proposal

A proposal can be submitted to update one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "update existing entries",
  "description": "update existing entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgUpdateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2",
      "decimals": "18",
      "denom": "mytoken2",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken2",
      "displaySymbol": "mytoken2",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken2",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken2",
      "ibcCounterpartyDenom": "mytoken2",
      "ibcCounterpartyChainId": "test"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Delete Entry using Gov Proposal

A proposal can be submitted to delete one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "delete entries",
  "description": "delete entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgDeleteEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```
