#!/bin/bash

# Function to calculate SHA-256 hash and format it
calculate_hash() {
    local CHANNEL_ID="$1"
    local BASE_DENOM="$2"
    local DENOM
    DENOM=$(echo -n "transfer/channel-$CHANNEL_ID/$BASE_DENOM" | shasum -a 256 | tr '[:lower:]' '[:upper:]' | sed 's/.$//')
    echo "$DENOM"
}

# Function to generate JSON for a given set of parameters
generate_json() {
    local BASE_DENOM="$1"
    local NAME="$2"
    local DECIMALS="$3"
    local CHANNEL_ID="$4"
    local EXTERNAL_CHANNEL_ID="$5"
    local DENOM
    DENOM=$(calculate_hash "$CHANNEL_ID" "$BASE_DENOM")

    cat <<EOF
    {
      "@type": "/elys.oracle.MsgAddAssetInfo",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "denom": "ibc/$DENOM",
      "display": "$NAME",
      "band_ticker": "$NAME",
      "elys_ticker": "$NAME",
      "decimal": "$DECIMALS"
    },
EOF
}

# Main script
if [ "$#" -lt 5 ]; then
    echo "Usage: $0 <base_denom1> <name1> <decimal1> <channel_id1> <external_channel_id1> [<base_denom2> <name2> <decimal2> <channel_id2> <external_channel_id2> ...]"
    exit 1
fi

cat <<EOF
{
  "title": "add denom entries",
  "summary": "add denom entries",
  "messages": [
EOF

# Loop through the arguments and generate JSON for each set of inputs
while [ "$#" -ge 5 ]; do
    BASE_DENOM="$1"
    NAME="$2"
    DECIMALS="$3"
    CHANNEL_ID="$4"
    EXTERNAL_CHANNEL_ID="$5"
    generate_json "$BASE_DENOM" "$NAME" "$DECIMALS" "$CHANNEL_ID" "$EXTERNAL_CHANNEL_ID"
    shift 5
done

cat <<EOF
  ],
  "deposit": "10000000uelys"
}
EOF