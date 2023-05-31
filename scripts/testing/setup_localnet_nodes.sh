#!/bin/bash

# Function to check if screen is installed
check_screen_installed() {
    if ! command -v screen &> /dev/null; then
        echo "Error: screen is not installed. Please install screen and try again."
        exit 1
    fi
}

# Function to check if any required environment variables are missing
check_required_variables() {
    local missing_variables=()
    for var in "${required_variables[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_variables+=("$var")
        fi
    done

    if [[ ${#missing_variables[@]} -gt 0 ]]; then
        echo "The following environment variables are missing or empty:"
        for var in "${missing_variables[@]}"; do
            echo "- $var"
        done
        exit 1
    fi

    echo "All required environment variables are set."
}

# Function to prompt user to stop existing screen sessions
prompt_to_stop_screens() {
    local existing_screens=()
    for folder in "${NODE_FOLDERS[@]}"; do
        screen -ls "$folder" | grep -q "$folder"
        if [ $? -eq 0 ]; then
            existing_screens+=("$folder")
        fi
    done

    if [[ ${#existing_screens[@]} -gt 0 ]]; then
        echo "The following screens already exist:"
        for screen_name in "${existing_screens[@]}"; do
            echo "- $screen_name"
        done

        read -p "Do you want to stop these screens? (y/n): " choice
        case $choice in
            [Yy]*)
                for screen_name in "${existing_screens[@]}"; do
                    screen -S "$screen_name" -X quit
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
}

# Function to prompt user to delete existing node folders
prompt_to_delete_folders() {
    local existing_folders=()
    for folder in "${NODE_FOLDERS[@]}"; do
        if [[ -d "/tmp/$folder" ]]; then
            existing_folders+=("$folder")
        fi
    done

    if [[ ${#existing_folders[@]} -gt 0 ]]; then
        echo "The following node folders already exist:"
        for folder in "${existing_folders[@]}"; do
            echo "- $folder"
        done

        read -p "Do you want to delete these folders? (y/n): " choice
        case $choice in
            [Yy]*)
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
}

# Function to download snapshot data
download_snapshot_data() {
    if [[ -f "${SNAPSHOT_PATH}" ]]; then
        echo "Snapshot file already exists. Skipping download..."
    else
        echo "Downloading testnet snapshot data..."
        wget -O "${SNAPSHOT_PATH}" --inet4-only "${SNAPSHOT_URL}"
    fi
}

# Function to initialize nodes
initialize_nodes() {
    for index in "${!NODE_FOLDERS[@]}"; do
        local folder="${NODE_FOLDERS[index]}"
        local node_rpc_port="${NODE_RPC_PORTS[index]}"
        echo "Initializing $folder node..."
        ${BINARY} init "$folder" --chain-id "${CHAIN_ID}" --home "/tmp/$folder" >/dev/null 2>&1
        echo "Configuring $folder node..."
        ${BINARY} config node "tcp://localhost:${node_rpc_port}" --home "/tmp/$folder" >/dev/null 2>&1
        ${BINARY} config keyring-backend test --home "/tmp/$folder" >/dev/null 2>&1
        ${BINARY} config broadcast-mode block --home "/tmp/$folder" >/dev/null 2>&1
    done
}

# Function to update files
update_files() {
    for index in "${!NODE_FOLDERS[@]}"; do
        local folder="${NODE_FOLDERS[index]}"
        local validator_address="${VALIDATOR_ADDRESSES[index]}"
        local validator_public_key="${VALIDATOR_PUBLIC_KEYS[index]}"
        local validator_private_key="${VALIDATOR_PRIVATE_KEYS[index]}"
        local node_private_key="${NODE_PRIVATE_KEYS[index]}"
        local peers=""
        for index2 in "${!NODE_FOLDERS[@]}"; do
            peers+="${NODE_IDS[index2]}@0.0.0.0:${NODE_P2P_PORTS[index2]},"
        done
        local node_rpc_port="${NODE_RPC_PORTS[index]}"
        local node_p2p_port="${NODE_P2P_PORTS[index]}"
        local node_grpc_port="${NODE_GRPC_PORTS[index]}"
        local node_grpc_web_port="${NODE_GRPC_WEB_PORTS[index]}"
        local node_pprof_port="${NODE_PPROF_PORTS[index]}"

        echo "Updating files for $folder node..."
        local priv_validator_key_path="/tmp/$folder/config/priv_validator_key.json"
        local node_key_path="/tmp/$folder/config/node_key.json"
        local config_path="/tmp/$folder/config/config.toml"
        local app_path="/tmp/$folder/config/app.toml"

        # Update priv_validator_key.json
        local priv_validator_key_content=$(cat <<EOF
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
        echo "$priv_validator_key_content" > "$priv_validator_key_path"

        # Update node_key.json
        local node_key_content=$(cat <<EOF
{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"${node_private_key}"}}
EOF
        )
        echo "$node_key_content" > "$node_key_path"

        # Update config file peers
        peers=${peers::-1} # Remove last comma
        sed -i "s/^persistent_peers =.*/persistent_peers = \"${peers}\"/" "$config_path"

        # Set config file allow_duplicate_ip property to true
        sed -i "s/^allow_duplicate_ip =.*/allow_duplicate_ip = true/" "$config_path"

        # Set config file allow_duplicate_ip property to true
        sed -i "s/^cors_allowed_origins =.*/cors_allowed_origins = [\"*\"]/" "$config_path"

        # Update rpc laddr in config with node port
        sed -i "s/^laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/127.0.0.1:${node_rpc_port}\"/" "$config_path"

        # Update p2p laddr in config with node port
        sed -i "s/^laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/127.0.0.1:${node_p2p_port}\"/" "$config_path"

        # Update p2p laddr in config with node port
        sed -i "s/^pprof_laddr =.*/pprof_laddr = \"localhost:${node_pprof_port}\"/" "$config_path"

        # Update grpc in app with node port
        sed -i "s/^address = \"0.0.0.0:9090\"/address = \"0.0.0.0:${node_grpc_port}\"/" "$app_path"

        # Update grpc web in app with node port
        sed -i "s/^address = \"0.0.0.0:9091\"/address = \"0.0.0.0:${node_grpc_web_port}\"/" "$app_path"
    done
}

# Function to import validator account
import_validator_account() {
    for index in "${!NODE_FOLDERS[@]}"; do
        local folder="${NODE_FOLDERS[index]}"
        local validator_private_key_path="/tmp/${folder}/validator_private_key.pem"
        local validator_private_key_salt="${VALIDATOR_PRIVATE_KEY_SALTS[index]}"
        local validator_private_key_first_line="${VALIDATOR_PRIVATE_KEY_FIRST_LINES[index]}"
        local validator_private_key_second_line="${VALIDATOR_PRIVATE_KEY_SECOND_LINES[index]}"
        local validator_private_key_third_line="${VALIDATOR_PRIVATE_KEY_THIRD_LINES[index]}"
        local validator_account_passphrase="${VALIDATOR_ACCOUNT_PASSPHRASES[index]}"

        echo "Importing validator account to $folder node..."

        # Update validator private key file
        local priv_validator_key_content=$(cat <<EOF
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

        echo "$priv_validator_key_content" > "$validator_private_key_path"
        ${BINARY} keys import validator "$validator_private_key_path" --keyring-backend test --home "/tmp/$folder" <<< "$validator_account_passphrase"
    done
}

# Function to extract and place data folder
extract_data_folder() {
    for folder in "${NODE_FOLDERS[@]}"; do
        echo "Extracting data folder for $folder node..."
        lz4 -c -d "${SNAPSHOT_PATH}" | tar -x -C "/tmp/$folder"
    done
}

# Function to add missing upgrade info to data folder
add_missing_upgrade_info() {
    for folder in "${NODE_FOLDERS[@]}"; do
        echo "Adding missing upgrade info to data folder for $folder node..."

        local upgrade_info_content=$(cat <<EOF
{
    "name": "v0.5.3",
    "time": "0001-01-01T00:00:00Z",
    "height": 477000,
    "info": "{\"binaries\":{\"linux/amd64\":\"https://github.com/elys-network/elys/releases/download/v0.5.3/elys._v0.5.3_linux_amd64.tar.gz?checksum=4d8824eb30e416420666ffc30e74b7d57c1f913b89e38f24aae889660d322c0f\"}}"
}
EOF
        )

        echo "${upgrade_info_content}" > "/tmp/$folder/data/upgrade-info.json"
    done
}

# Function to start nodes in separate screen sessions
start_nodes_in_screens() {
    for folder in "${NODE_FOLDERS[@]}"; do
        echo "Starting $folder node in a screen session..."

        screen -dmS "$folder" "${BINARY}" start --home "/tmp/$folder"
        sleep 1
    done
}

# Function to display all screen sessions
display_screen_sessions() {
    screen -ls
    echo "You can attach to a screen with the following command:"
    echo "screen -r <screen_name>"
    echo "Example: screen -r node0"
    echo "You can detach from a screen by pressing Ctrl+A and then Ctrl+D"
}

# Function to submit software upgrade proposal
submit_upgrade_proposal() {
    local first_folder="${NODE_FOLDERS[0]}"

    ${BINARY} tx gov submit-legacy-proposal software-upgrade \
        v0.6.1 \
        --deposit=10000000uelys \
        --upgrade-height=770750 \
        --title="v0.6.1" \
        --description="v0.6.1" \
        --no-validate \
        --from=validator \
        --fees=100000uelys \
        --gas=auto \
        --home="/tmp/${first_folder}" \
        -y

    for folder in "${NODE_FOLDERS[@]}"; do
        ${BINARY} tx gov vote \
            11 \
            yes \
            --from=validator \
            --fees=100000uelys \
            --home="/tmp/${folder}" \
            -y
    done
}

# Function to prompt user to submit software upgrade proposal
prompt_to_submit_software_upgrade() {
    read -p "Do you want to submit a software upgrade proposal? (y/n): " choice
    case $choice in
        [Yy]*)
            submit_upgrade_proposal
            ;;
        [Nn]*)
            echo "No software upgrade proposal submitted."
            ;;
        *)
            echo "Invalid choice. No software upgrade proposal submitted."
            ;;
    esac
}

# Main script execution

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
    "CHAIN_ID"
    "SNAPSHOT_URL"
    "SNAPSHOT_PATH"
)

# Check if screen is installed
check_screen_installed

# Check if any required variable is missing
check_required_variables

# Check if any of the named screens already exist and prompt to stop them
prompt_to_stop_screens

# Prompt to delete existing node folders
prompt_to_delete_folders

# Continue with the script
echo "Continuing with node setup..."

# Download snapshot data
download_snapshot_data

# Initialize nodes
initialize_nodes

# Update files
update_files

# Import validator account
import_validator_account

# Extract and place data folder
extract_data_folder

# Add missing upgrade info to data folder
add_missing_upgrade_info

# Start nodes in separate screen sessions
start_nodes_in_screens

# Display all screen sessions
display_screen_sessions

# Prompt user to submit a software upgrade proposal
prompt_to_submit_software_upgrade

# Check if any of the named screens exist and prompt to stop them
prompt_to_stop_screens