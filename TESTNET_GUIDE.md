# FutchainEVM Local Testnet Guide

This guide helps you set up and test your FutchainEVM blockchain locally with pre-funded test accounts.

## Quick Start

### 1. Build and Start the Testnet

```bash
# Build and start with fresh data
./run_testnet.sh -y --show-accounts

# Or start without rebuilding
./run_testnet.sh -n --show-accounts
```

### 2. Test Your Blockchain

```bash
# Check node status
./test_commands.sh status

# View account balances
./test_commands.sh balances

# Test futchain module
./test_commands.sh futchain

# Send test transaction
./test_commands.sh send alice bob 5000000

# Monitor logs
./test_commands.sh logs
```

## Test Accounts

The testnet comes with 7 pre-funded accounts:

| Account   | Purpose                    | Balance    | Key Name  |
|-----------|----------------------------|------------|-----------|
| validator | Main validator account     | 100M tokens| mykey     |
| alice     | Test user 1               | 1M tokens  | alice     |
| bob       | Test user 2               | 1M tokens  | bob       |
| charlie   | Test user 3               | 1M tokens  | charlie   |
| diana     | Test user 4               | 1M tokens  | diana     |
| eve       | Test user 5               | 500K tokens| eve       |
| frank     | Test user 6               | 500K tokens| frank     |

## Available Endpoints

When the node is running, you can access:

- **Cosmos RPC**: http://localhost:26657
- **Cosmos REST API**: http://localhost:1317
- **Ethereum JSON-RPC**: http://localhost:8545
- **WebSocket**: ws://localhost:8546

## Script Options

### run_testnet.sh Options

```bash
./run_testnet.sh [OPTIONS]

Options:
  -y                 Overwrite existing chain data without prompt
  -n                 Keep existing chain data without prompt
  --no-install       Skip building/installing the binary
  --debug            Build with debug options
  --show-accounts    Display test account information
  --help, -h         Show help message
```

### test_commands.sh Commands

```bash
./test_commands.sh [COMMAND] [OPTIONS]

Commands:
  status              Show node status and endpoints
  balances            Show all account balances
  futchain            Test futchain module queries
  send [from] [to] [amount]  Send test transaction
  gov                 Submit test governance proposal
  logs                Monitor node logs
  help                Show help message
```

## Testing Your Futchain Module

### Query Module Parameters

```bash
futchaind query futchain params --home ~/.futchaind
```

### Monitor Data Fetching

Your futchain module fetches sports data every 10 blocks by default. Watch the logs to see it in action:

```bash
./test_commands.sh logs
```

Look for log entries like:
```
[INFO] fetching data block height=20 fetch modulo=10
[INFO] detected a new league league="Premier League" id=39 event=new_league
[INFO] detected a new match match=123 event=new_match
```

### Test Governance

Submit a proposal to change futchain parameters:

```bash
./test_commands.sh gov
```

Then vote on the proposal:

```bash
# List proposals
futchaind query gov proposals --home ~/.futchaind

# Vote on proposal (replace 1 with actual proposal ID)
futchaind tx gov vote 1 yes --from validator --keyring-backend test --home ~/.futchaind --chain-id cosmos_262144-1 --yes
```

## Common Commands

### Check Account Balance

```bash
futchaind query bank balances $(futchaind keys show alice -a --keyring-backend test --home ~/.futchaind) --home ~/.futchaind
```

### Send Tokens

```bash
futchaind tx bank send alice $(futchaind keys show bob -a --keyring-backend test --home ~/.futchaind) 1000000utest \
  --keyring-backend test --home ~/.futchaind --chain-id cosmos_262144-1 --yes
```

### Query Transaction

```bash
futchaind query tx [TX_HASH] --home ~/.futchaind
```

### View Validator Info

```bash
futchaind query staking validators --home ~/.futchaind
```

## Troubleshooting

### Node Won't Start

1. Check if ports are already in use:
   ```bash
   lsof -i :26657  # Cosmos RPC
   lsof -i :8545   # Ethereum RPC
   ```

2. Clean up and restart:
   ```bash
   ./run_testnet.sh -y
   ```

### Can't Connect to Node

1. Verify node is running:
   ```bash
   curl http://localhost:26657/status
   ```

2. Check logs:
   ```bash
   tail -f ~/.futchaind/logs/futchaind.log
   ```

### Account Not Found

Make sure you're using the correct keyring backend and home directory:

```bash
futchaind keys list --keyring-backend test --home ~/.futchaind
```

## Development Tips

1. **Use test keyring**: The scripts use `--keyring-backend test` for convenience. Never use this in production!

2. **Monitor logs**: Keep an eye on logs to see your futchain module fetching data.

3. **Fast governance**: Governance proposals have short voting periods (30s) for quick testing.

4. **Reset anytime**: Use `./run_testnet.sh -y` to completely reset the blockchain state.

5. **Multiple terminals**: Run the node in one terminal and use `./test_commands.sh` in another.

## Next Steps

- Deploy smart contracts using the Ethereum JSON-RPC endpoint
- Test IBC transfers with another chain
- Experiment with ERC20 token conversions
- Build frontend applications using the REST API
- Test your futchain module's data fetching and storage

Happy testing! 🚀
