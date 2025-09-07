#!/bin/bash

# FutchainEVM Local Testnet Script
# This script sets up and runs a local testnet with pre-funded test accounts
# for development and testing purposes.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHAINID="${CHAIN_ID:-cosmos_262144-1}"
MONIKER="futchain-testnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
HOMEDIR="$HOME/.futchaind"
BASEFEE=10000000
DENOM="utest"

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Validate dependencies
command -v jq >/dev/null 2>&1 || {
    print_error "jq not installed. Install it from: https://stedolan.github.io/jq/download/"
    exit 1
}

# Parse input flags
install=true
overwrite=""
BUILD_FOR_DEBUG=false
SHOW_ACCOUNTS=false

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
    -y)
        print_status "Flag -y passed -> Overwriting the previous chain data."
        overwrite="y"
        shift
        ;;
    -n)
        print_status "Flag -n passed -> Not overwriting the previous chain data."
        overwrite="n"
        shift
        ;;
    --no-install)
        print_status "Flag --no-install passed -> Skipping installation of the futchaind binary."
        install=false
        shift
        ;;
    --debug)
        print_status "Flag --debug passed -> Building with debug options."
        BUILD_FOR_DEBUG=true
        shift
        ;;
    --show-accounts)
        print_status "Flag --show-accounts passed -> Will display account information."
        SHOW_ACCOUNTS=true
        shift
        ;;
    --help|-h)
        echo "FutchainEVM Local Testnet Script"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  -y                 Overwrite existing chain data without prompt"
        echo "  -n                 Keep existing chain data without prompt"
        echo "  --no-install       Skip building/installing the binary"
        echo "  --debug            Build with debug options"
        echo "  --show-accounts    Display test account information"
        echo "  --help, -h         Show this help message"
        echo ""
        echo "Test Accounts:"
        echo "  validator  - Main validator account (100M tokens)"
        echo "  alice      - Test user 1 (1M tokens)"
        echo "  bob        - Test user 2 (1M tokens)"
        echo "  charlie    - Test user 3 (1M tokens)"
        echo "  diana      - Test user 4 (1M tokens)"
        echo "  eve        - Test user 5 (500K tokens)"
        echo "  frank      - Test user 6 (500K tokens)"
        echo ""
        exit 0
        ;;
    *)
        print_error "Unknown flag: $key"
        echo "Use --help for usage information."
        exit 1
        ;;
    esac
done

print_header "FutchainEVM Local Testnet Setup"

# Check if binary exists (unless we're building it)
if [[ $install == false ]]; then
    command -v futchaind >/dev/null 2>&1 || {
        print_error "futchaind binary not found. Run with default options to build it, or run 'make futchaind' first."
        exit 1
    }
fi

# Build binary if requested
if [[ $install == true ]]; then
    print_status "Building futchaind binary..."
    if [[ $BUILD_FOR_DEBUG == true ]]; then
        CGO_ENABLED="1" go build -gcflags "all=-N -l" -o futchaind ./futchaind
    else
        make futchaind
    fi
    print_status "Binary built successfully!"
fi

# Check for existing configuration
if [[ $overwrite = "" ]]; then
    if [ -d "$HOMEDIR" ]; then
        print_warning "Existing configuration found at '$HOMEDIR'"
        echo "Overwrite the existing configuration and start fresh? [y/n]"
        read -r overwrite
    else
        overwrite="y"
    fi
fi

# Test accounts with mnemonics (compatible with Bash 3.x)
# Format: account_name:key_name:mnemonic:balance

ACCOUNT_DATA=(
    "validator:mykey:gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat:100000000000000000000000000"
    "alice:alice:copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom:1000000000000000000000000"
    "bob:bob:maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual:1000000000000000000000000"
    "charlie:charlie:will wear settle write dance topic tape sea glory hotel oppose rebel client problem era video gossip glide during yard balance cancel file rose:1000000000000000000000000"
    "diana:diana:doll midnight silk carpet brush boring pluck office gown inquiry duck chief aim exit gain never tennis crime fragile ship cloud surface exotic patch:1000000000000000000000000"
)

# Setup function
setup_chain() {
    print_header "Setting up fresh blockchain"
    
    # Remove existing data
    rm -rf "$HOMEDIR"
    
    # Set client config
    futchaind config set client chain-id "$CHAINID" --home "$HOMEDIR"
    futchaind config set client keyring-backend "$KEYRING" --home "$HOMEDIR"
    
    print_status "Importing test accounts..."
    
    # Import all accounts
    for account_info in "${ACCOUNT_DATA[@]}"; do
        IFS=':' read -r account key_name mnemonic balance <<< "$account_info"
        echo "$mnemonic" | futchaind keys add "$key_name" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR" >/dev/null 2>&1
        print_status "✓ Imported account: $account ($key_name)"
    done
    
    # Initialize chain
    print_status "Initializing blockchain..."
    futchaind init $MONIKER -o --chain-id "$CHAINID" --home "$HOMEDIR" >/dev/null 2>&1
    
    # Configure genesis parameters
    print_status "Configuring genesis parameters..."
    
    # Set denomination
    jq --arg denom "$DENOM" '.app_state["staking"]["params"]["bond_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state["gov"]["params"]["min_deposit"][0]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state["evm"]["params"]["evm_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state["mint"]["params"]["mint_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure futchain module parameters
    jq '.app_state["futchain"]["params"]["fetch_modulo"]="3"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["futchain"]["params"]["timezone"]="Europe/Istanbul"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Enable EVM precompiles
    jq '.app_state["evm"]["params"]["active_static_precompiles"]=["0x0000000000000000000000000000000000000100","0x0000000000000000000000000000000000000400","0x0000000000000000000000000000000000000800","0x0000000000000000000000000000000000000801","0x0000000000000000000000000000000000000802","0x0000000000000000000000000000000000000803","0x0000000000000000000000000000000000000804","0x0000000000000000000000000000000000000805"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Enable native denomination as token pair
    jq '.app_state.erc20.params.native_precompiles=["0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state.erc20.token_pairs=[{contract_owner:1,erc20_address:"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",denom:$denom,enabled:true}]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Set gas limit
    jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure governance for faster testing
    jq '.app_state["gov"]["deposit_params"]["max_deposit_period"]="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["gov"]["voting_params"]["voting_period"]="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["gov"]["params"]["max_deposit_period"]="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["gov"]["params"]["voting_period"]="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state["gov"]["params"]["expedited_voting_period"]="15s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure node settings
    print_status "Configuring node settings..."
    
    # Enable APIs and metrics
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/prometheus = false/prometheus = true/' "$CONFIG"
        sed -i '' 's/prometheus-retention-time = 0/prometheus-retention-time = 1000000000000/g' "$APP_TOML"
        sed -i '' 's/enabled = false/enabled = true/g' "$APP_TOML"
        sed -i '' 's/enable = false/enable = true/g' "$APP_TOML"
    else
        sed -i 's/prometheus = false/prometheus = true/' "$CONFIG"
        sed -i 's/prometheus-retention-time = "0"/prometheus-retention-time = "1000000000000"/g' "$APP_TOML"
        sed -i 's/enabled = false/enabled = true/g' "$APP_TOML"
        sed -i 's/enable = false/enable = true/g' "$APP_TOML"
    fi
    
    # Set pruning settings
    sed -i.bak 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
    sed -i.bak 's/pruning-keep-recent = "0"/pruning-keep-recent = "2"/g' "$APP_TOML"
    sed -i.bak 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"
    
    # Add genesis accounts
    print_status "Adding genesis accounts..."
    for account_info in "${ACCOUNT_DATA[@]}"; do
        IFS=':' read -r account key_name mnemonic balance <<< "$account_info"
        balance_with_denom="${balance}$DENOM"
        futchaind genesis add-genesis-account "$key_name" "$balance_with_denom" --keyring-backend "$KEYRING" --home "$HOMEDIR" >/dev/null 2>&1
        print_status "✓ Added genesis account: $account ($balance_with_denom)"
    done
    
    # Create validator
    print_status "Creating validator..."
    validator_key="mykey"  # First account is always validator
    futchaind genesis gentx "$validator_key" "1000000000000000000000$DENOM" --gas-prices "${BASEFEE}$DENOM" --keyring-backend "$KEYRING" --chain-id "$CHAINID" --home "$HOMEDIR" >/dev/null 2>&1
    
    # Collect genesis transactions
    futchaind genesis collect-gentxs --home "$HOMEDIR" >/dev/null 2>&1
    
    # Validate genesis
    print_status "Validating genesis..."
    futchaind genesis validate-genesis --home "$HOMEDIR" >/dev/null 2>&1
    
    print_status "✓ Blockchain setup complete!"
}

# Display account information
show_account_info() {
    print_header "Test Account Information"
    
    echo -e "${BLUE}Chain ID:${NC} $CHAINID"
    echo -e "${BLUE}Denomination:${NC} $DENOM"
    echo -e "${BLUE}Home Directory:${NC} $HOMEDIR"
    echo ""
    
    printf "%-12s %-45s %-20s %s\n" "Account" "Address" "Balance" "Purpose"
    printf "%-12s %-45s %-20s %s\n" "--------" "-------" "-------" "-------"
    
    for account_info in "${ACCOUNT_DATA[@]}"; do
        IFS=':' read -r account key_name mnemonic balance_raw <<< "$account_info"
        address=$(futchaind keys show "$key_name" -a --keyring-backend "$KEYRING" --home "$HOMEDIR" 2>/dev/null || echo "N/A")
        
        # Convert balance to human readable format
        if [[ ${#balance_raw} -gt 18 ]]; then
            balance_int=${balance_raw:0:-18}
            balance_dec=${balance_raw: -18}
            balance="${balance_int}.${balance_dec:0:6}M"
        else
            balance="$balance_raw"
        fi
        
        case $account in
            validator) purpose="Validator & Main Account" ;;
            alice|bob|charlie|diana) purpose="Test User (1M tokens)" ;;
            eve|frank) purpose="Test User (500K tokens)" ;;
        esac
        
        printf "%-12s %-45s %-20s %s\n" "$account" "$address" "$balance" "$purpose"
    done
    
    echo ""
    print_status "Use these accounts for testing transactions and smart contracts"
    print_status "Keyring backend: $KEYRING (for development only!)"
}

# Main execution
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
    setup_chain
fi

if [[ $SHOW_ACCOUNTS == true ]]; then
    show_account_info
    echo ""
fi

# Display useful commands
print_header "Useful Commands"
echo "Query futchain module:"
echo "  futchaind query futchain params --home $HOMEDIR"
echo ""
echo "Check account balance:"
echo "  futchaind query bank balances \$(futchaind keys show alice -a --keyring-backend test --home $HOMEDIR) --home $HOMEDIR"
echo ""
echo "Send tokens:"
echo "  futchaind tx bank send alice \$(futchaind keys show bob -a --keyring-backend test --home $HOMEDIR) 1000000$DENOM --keyring-backend test --home $HOMEDIR --chain-id $CHAINID --yes"
echo ""
echo "View logs:"
echo "  tail -f $HOMEDIR/logs/futchaind.log"
echo ""

# Start the node
print_header "Starting FutchainEVM Node"
print_status "Chain ID: $CHAINID"
print_status "Home: $HOMEDIR"
print_status "Log Level: $LOGLEVEL"
print_status "Minimum Gas Price: 0.0001$DENOM"
print_warning "Press Ctrl+C to stop the node"
echo ""

# Create logs directory
mkdir -p "$HOMEDIR/logs"

# Start node with logging
futchaind start \
    --log_level $LOGLEVEL \
    --minimum-gas-prices="0.0001$DENOM" \
    --home "$HOMEDIR" \
    --json-rpc.api eth,txpool,personal,net,debug,web3 \
    --chain-id "$CHAINID"
