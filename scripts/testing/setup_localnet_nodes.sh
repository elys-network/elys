#!/bin/bash

# Load environment variables
source scripts/testing/variables.sh

# List of environment variables to check
required_variables=(
    "NODE_FOLDERS"
    "NODE_RPC_PORTS"
    "NODE_GRPC_PORTS"
    "NODE_GRPC_WEB_PORTS"
    "NODE_P2P_PORTS"
    "NODE_PPROF_PORTS"
    "NODE_IDS"
    "VALIDATOR_ADDRESSES"
    "VALIDATOR_PUBLIC_KEYS"
    "VALIDATOR_PRIVATE_KEYS"
    "NODE_PRIVATE_KEYS"
    "VALIDATOR_PRIVATE_KEY_SALTS"
    "VALIDATOR_PRIVATE_KEY_FIRST_LINES"
    "VALIDATOR_PRIVATE_KEY_SECOND_LINES"
    "VALIDATOR_PRIVATE_KEY_THIRD_LINES"
    "VALIDATOR_ACCOUNT_PASSPHRASES"
    "BINARY"
    "SNAPSHOT_URL"
    "SNAPSHOT_PATH"
)

# Check if any required variable is missing
missing_variables=()
for var in "${required_variables[@]}"; do
    if [[ -z "${!var}" ]]; then
        missing_variables+=("$var")
    fi
done

# If any variable is missing, display error and exit
if [[ ${#missing_variables[@]} -gt 0 ]]; then
    echo "The following environment variables are missing or empty:"
    for var in "${missing_variables[@]}"; do
        echo "- $var"
    done
    exit 1
fi

# All required variables are present
echo "All required environment variables are set."

# Check if node folders already exist
existing_folders=()
for folder in "${NODE_FOLDERS[@]}"; do
    if [[ -d "/tmp/$folder" ]]; then
        existing_folders+=("$folder")
    fi
done

# Prompt user for deletion confirmation
if [[ ${#existing_folders[@]} -gt 0 ]]; then
    echo "The following node folders already exist:"
    for folder in "${existing_folders[@]}"; do
        echo "- $folder"
    done

    read -p "Do you want to delete these folders? (y/n): " choice
    case $choice in
        [Yy]*)
            # Delete existing node folders
            for folder in "${existing_folders[@]}"; do
                echo "Deleting folder: $folder"
                rm -rf "/tmp/$folder"
            done
            ;;
        [Nn]*)
            echo "Exiting..."
            exit 0
            ;;
        *)
            echo "Invalid choice. Exiting..."
            exit 1
            ;;
    esac
fi

# Continue with the script
echo "Continuing with node setup..."
# Rest of your script here

# Step 0: Check if snapshot file already exist
if [[ -f "${SNAPSHOT_PATH}" ]]; then
    echo "Snapshot file already exists. Skipping download..."
else
    echo "Downloading testnet snapshot data..."
    wget -O ${SNAPSHOT_PATH} --inet4-only ${SNAPSHOT_URL}
fi

# Step 1: Initialize nodes
for folder in "${NODE_FOLDERS[@]}"; do
    echo "Initializing $folder node..."
    ${BINARY} init $folder --home /tmp/$folder >/dev/null 2>&1
done

# Step 2: Update files
for index in "${!NODE_FOLDERS[@]}"; do
    folder="${NODE_FOLDERS[index]}"
    validator_address="${VALIDATOR_ADDRESSES[index]}"
    validator_public_key="${VALIDATOR_PUBLIC_KEYS[index]}"
    validator_private_key="${VALIDATOR_PRIVATE_KEYS[index]}"
    node_private_key="${NODE_PRIVATE_KEYS[index]}"
    peers=""
    for index2 in "${!NODE_FOLDERS[@]}"; do
        peers+="${NODE_IDS[index2]}@0.0.0.0:${NODE_P2P_PORTS[index2]},"
    done
    node_rpc_port="${NODE_RPC_PORTS[index]}"
    node_p2p_port="${NODE_P2P_PORTS[index]}"
    node_grpc_port="${NODE_GRPC_PORTS[index]}"
    node_grpc_web_port="${NODE_GRPC_WEB_PORTS[index]}"
    node_pprof_port="${NODE_PPROF_PORTS[index]}"

    echo "Updating files for $folder node..."
    priv_validator_key_path="/tmp/$folder/config/priv_validator_key.json"
    node_key_path="/tmp/$folder/config/node_key.json"
    config_path="/tmp/$folder/config/config.toml"
    app_path="/tmp/$folder/config/app.toml"

    # Update priv_validator_key.json
    priv_validator_key_content=$(cat <<EOF
{
  "address": "${validator_address}",
  "pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "${validator_public_key}"
  },
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "${validator_private_key}"
  }
}
EOF
)
    echo "$priv_validator_key_content" > $priv_validator_key_path

    # Update node_key.json
    node_key_content=$(cat <<EOF
{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"${node_private_key}"}}
EOF
)
    echo "$node_key_content" > $node_key_path

    # Update config file peers
    peers=${peers::-1} # Remove last comma
    sed -i "s/^persistent_peers =.*/persistent_peers = \"${peers}\"/" $config_path

    # Set config file allow_duplicate_ip property to true
    sed -i "s/^allow_duplicate_ip =.*/allow_duplicate_ip = true/" $config_path

    # Set config file allow_duplicate_ip property to true
    sed -i "s/^cors_allowed_origins =.*/cors_allowed_origins = [\"*\"]/" $config_path

    # Update rpc laddr in config with node port
    sed -i "s/^laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/127.0.0.1:${node_rpc_port}\"/" $config_path

    # Update p2p laddr in config with node port
    sed -i "s/^laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/127.0.0.1:${node_p2p_port}\"/" $config_path

    # Update p2p laddr in config with node port
    sed -i "s/^pprof_laddr =.*/pprof_laddr = \"localhost:${node_pprof_port}\"/" $config_path

    # Update grpc in app with node port
    sed -i "s/^address = \"0.0.0.0:9090\"/address = \"0.0.0.0:${node_grpc_port}\"/" $app_path

    # Update grpc web in app with node port
    sed -i "s/^address = \"0.0.0.0:9091\"/address = \"0.0.0.0:${node_grpc_web_port}\"/" $app_path
done

# Step 3: Import validator account
for index in "${!NODE_FOLDERS[@]}"; do
    folder="${NODE_FOLDERS[index]}"
    validator_private_key_path="/tmp/${folder}/validator_private_key.pem"
    validator_private_key_salt="${VALIDATOR_PRIVATE_KEY_SALTS[index]}"
    validator_private_key_first_line="${VALIDATOR_PRIVATE_KEY_FIRST_LINES[index]}"
    validator_private_key_second_line="${VALIDATOR_PRIVATE_KEY_SECOND_LINES[index]}"
    validator_private_key_third_line="${VALIDATOR_PRIVATE_KEY_THIRD_LINES[index]}"
    validator_account_passphrase="${VALIDATOR_ACCOUNT_PASSPHRASES[index]}"

    echo "Importing validator account to $folder node..."

    # Update validator private key file
    priv_validator_key_content=$(cat <<EOF
-----BEGIN TENDERMINT PRIVATE KEY-----
kdf: bcrypt
salt: ${validator_private_key_salt}
type: secp256k1

${validator_private_key_first_line}
${validator_private_key_second_line}
${validator_private_key_third_line}
-----END TENDERMINT PRIVATE KEY-----
EOF
)

    echo "$priv_validator_key_content" > $validator_private_key_path
    ${BINARY} keys import validator $validator_private_key_path --keyring-backend test --home /tmp/$folder <<< "$validator_account_passphrase"
done

# Step 4: Extract and place data folder
for folder in "${NODE_FOLDERS[@]}"; do
    echo "Extracting data folder for $folder node..."
    lz4 -c -d ${SNAPSHOT_PATH} | tar -x -C /tmp/$folder
done

# # Step 5: Start nodes
# for folder in "${NODE_FOLDERS[@]}"; do
#     echo "Starting $folder node..."
#     ${BINARY} start --home /tmp/$folder
# done
