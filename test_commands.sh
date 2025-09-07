#!/bin/bash

# FutchainEVM Test Commands Script
# This script provides common commands for testing your blockchain

set -e

# Configuration
CHAINID="${CHAIN_ID:-cosmos_262144-1}"
KEYRING="test"
HOMEDIR="$HOME/.futchaind"
DENOM="utest"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_command() {
    echo -e "${YELLOW}Command:${NC} $1"
}

# Check if node is running
check_node() {
    if ! curl -s http://localhost:26657/status >/dev/null 2>&1; then
        echo "❌ Node is not running. Start it with: ./run_testnet.sh"
        exit 1
    fi
    echo "✅ Node is running"
}

# Show account balances
show_balances() {
    print_header "Account Balances"
    
    accounts=("validator" "alice" "bob" "charlie" "diana" "eve" "frank")
    
    for account in "${accounts[@]}"; do
        address=$(futchaind keys show "$account" -a --keyring-backend "$KEYRING" --home "$HOMEDIR" 2>/dev/null || echo "N/A")
        if [[ "$address" != "N/A" ]]; then
            balance=$(futchaind query bank balances "$address" --home "$HOMEDIR" -o json 2>/dev/null | jq -r ".balances[0].amount // \"0\"")
            # Convert to human readable
            if [[ ${#balance} -gt 18 ]]; then
                balance_int=${balance:0:-18}
                balance_dec=${balance: -18}
                readable_balance="${balance_int}.${balance_dec:0:6}M"
            else
                readable_balance="$balance"
            fi
            printf "%-12s %-45s %s %s\n" "$account" "$address" "$readable_balance" "$DENOM"
        fi
    done
}

# Test futchain module
test_futchain() {
    print_header "Testing Futchain Module"
    
    print_status "Querying futchain parameters..."
    print_command "futchaind query futchain params --home $HOMEDIR"
    futchaind query futchain params --home "$HOMEDIR"
    echo ""
    
    print_status "Checking current block height..."
    height=$(futchaind query block --home "$HOMEDIR" -o json | jq -r ".block.header.height")
    print_status "Current block height: $height"
    echo ""
    
    print_status "Note: Futchain module fetches data every 10 blocks by default"
    print_status "Watch the logs to see data fetching in action!"
}

# Send test transaction
send_test_tx() {
    local from_account=${1:-alice}
    local to_account=${2:-bob}
    local amount=${3:-1000000}
    
    print_header "Sending Test Transaction"
    
    from_addr=$(futchaind keys show "$from_account" -a --keyring-backend "$KEYRING" --home "$HOMEDIR")
    to_addr=$(futchaind keys show "$to_account" -a --keyring-backend "$KEYRING" --home "$HOMEDIR")
    
    print_status "From: $from_account ($from_addr)"
    print_status "To: $to_account ($to_addr)"
    print_status "Amount: $amount$DENOM"
    
    print_command "futchaind tx bank send $from_account $to_addr ${amount}$DENOM --keyring-backend $KEYRING --home $HOMEDIR --chain-id $CHAINID --yes"
    
    futchaind tx bank send "$from_account" "$to_addr" "${amount}$DENOM" \
        --keyring-backend "$KEYRING" \
        --home "$HOMEDIR" \
        --chain-id "$CHAINID" \
        --yes \
        --gas auto \
        --gas-adjustment 1.5
}

# Test governance proposal
test_governance() {
    print_header "Testing Governance"
    
    print_status "Creating a parameter change proposal..."
    
    # Create proposal JSON
    cat > /tmp/proposal.json <<EOF
{
  "messages": [
    {
      "@type": "/cosmos.gov.v1.MsgUpdateParams",
      "authority": "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
      "params": {
        "fetch_modulo": "5",
        "timezone": "UTC"
      }
    }
  ],
  "metadata": "Test proposal to change futchain parameters",
  "deposit": "10000000$DENOM",
  "title": "Update Futchain Parameters",
  "summary": "This proposal updates the fetch_modulo parameter from 10 to 5 blocks"
}
EOF
    
    print_command "futchaind tx gov submit-proposal /tmp/proposal.json --from validator --keyring-backend $KEYRING --home $HOMEDIR --chain-id $CHAINID --yes"
    
    futchaind tx gov submit-proposal /tmp/proposal.json \
        --from validator \
        --keyring-backend "$KEYRING" \
        --home "$HOMEDIR" \
        --chain-id "$CHAINID" \
        --yes \
        --gas auto \
        --gas-adjustment 1.5
    
    echo ""
    print_status "Proposal submitted! Check with: futchaind query gov proposals --home $HOMEDIR"
    
    # Clean up
    rm -f /tmp/proposal.json
}

# Monitor logs
monitor_logs() {
    print_header "Monitoring Node Logs"
    print_status "Press Ctrl+C to stop monitoring"
    echo ""
    
    if [[ -f "$HOMEDIR/logs/futchaind.log" ]]; then
        tail -f "$HOMEDIR/logs/futchaind.log"
    else
        print_status "Log file not found. Node might not be running with logging enabled."
        print_status "Try: journalctl -u futchaind -f (if running as systemd service)"
    fi
}

# Show node status
show_status() {
    print_header "Node Status"
    
    if curl -s http://localhost:26657/status >/dev/null 2>&1; then
        status=$(curl -s http://localhost:26657/status)
        
        network=$(echo "$status" | jq -r ".result.node_info.network")
        latest_height=$(echo "$status" | jq -r ".result.sync_info.latest_block_height")
        latest_time=$(echo "$status" | jq -r ".result.sync_info.latest_block_time")
        catching_up=$(echo "$status" | jq -r ".result.sync_info.catching_up")
        
        echo "Network: $network"
        echo "Latest Height: $latest_height"
        echo "Latest Block Time: $latest_time"
        echo "Catching Up: $catching_up"
        echo ""
        
        print_status "RPC Endpoints:"
        echo "  Cosmos RPC: http://localhost:26657"
        echo "  Cosmos REST: http://localhost:1317"
        echo "  Ethereum JSON-RPC: http://localhost:8545"
        echo "  WebSocket: ws://localhost:8546"
    else
        echo "❌ Node is not running or not accessible"
    fi
}

# Show help
show_help() {
    echo "FutchainEVM Test Commands"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  status              Show node status and endpoints"
    echo "  balances            Show all account balances"
    echo "  futchain            Test futchain module queries"
    echo "  send [from] [to] [amount]  Send test transaction (default: alice to bob, 1M tokens)"
    echo "  gov                 Submit test governance proposal"
    echo "  logs                Monitor node logs"
    echo "  help                Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status"
    echo "  $0 balances"
    echo "  $0 send alice bob 5000000"
    echo "  $0 futchain"
    echo "  $0 gov"
    echo "  $0 logs"
}

# Main execution
case "${1:-help}" in
    "status")
        check_node
        show_status
        ;;
    "balances")
        check_node
        show_balances
        ;;
    "futchain")
        check_node
        test_futchain
        ;;
    "send")
        check_node
        send_test_tx "$2" "$3" "$4"
        ;;
    "gov")
        check_node
        test_governance
        ;;
    "logs")
        monitor_logs
        ;;
    "help"|*)
        show_help
        ;;
esac
