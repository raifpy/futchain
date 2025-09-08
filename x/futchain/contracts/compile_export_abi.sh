#!/bin/bash

# Futchain Contract Compilation and ABI Export Script
# This script compiles the Solidity contracts and exports the ABI for Go embedding

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACTS_DIR="$SCRIPT_DIR"

echo -e "${BLUE}🔨 Futchain Contract Compilation Script${NC}"
echo -e "${BLUE}=======================================${NC}"
echo ""

# Check if solc is installed
if ! command -v solc &> /dev/null; then
    echo -e "${RED}❌ solc (Solidity compiler) is not installed${NC}"
    echo -e "${YELLOW}💡 Install it with:${NC}"
    echo "   npm install -g solc"
    echo "   # or"
    echo "   brew install solidity"
    exit 1
fi

echo -e "${GREEN}✅ Found solc version: $(solc --version | head -n1)${NC}"
echo ""

# Check if base.sol exists
if [ ! -f "$CONTRACTS_DIR/base.sol" ]; then
    echo -e "${RED}❌ base.sol not found in $CONTRACTS_DIR${NC}"
    exit 1
fi

echo -e "${BLUE}📁 Working directory: $CONTRACTS_DIR${NC}"
echo -e "${BLUE}📄 Compiling: base.sol${NC}"
echo ""

# Compile the contract and extract ABI
echo -e "${YELLOW}🔧 Compiling Solidity contract...${NC}"

# Create output directory if it doesn't exist
mkdir -p "$CONTRACTS_DIR/build"

# Compile and extract ABI for the FutI interface
solc --abi --optimize --output-dir "$CONTRACTS_DIR/build" "$CONTRACTS_DIR/base.sol" --overwrite

# Check if compilation was successful
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Compilation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Compilation successful${NC}"

# Find the FutI ABI file (it should be generated)
FUTI_ABI_FILE="$CONTRACTS_DIR/build/FutI.abi"

if [ -f "$FUTI_ABI_FILE" ]; then
    echo -e "${YELLOW}📋 Extracting FutI interface ABI...${NC}"
    
    # Copy the ABI to the contracts directory for Go embedding
    cp "$FUTI_ABI_FILE" "$CONTRACTS_DIR/abi.json"
    
    echo -e "${GREEN}✅ ABI exported to:${NC}"
    echo -e "   📄 $CONTRACTS_DIR/abi.json"
else
    echo -e "${YELLOW}⚠️  FutI.abi not found, trying to extract from base.sol manually...${NC}"
    
    # Alternative: Extract ABI using solc directly with interface filter
    solc --abi --optimize "$CONTRACTS_DIR/base.sol" 2>/dev/null | \
    awk '/======= .*FutI =======/{flag=1; next} /======= /{flag=0} flag && /^\[/{print; exit}' > "$CONTRACTS_DIR/abi.json"
    
    if [ -s "$CONTRACTS_DIR/abi.json" ]; then
        echo -e "${GREEN}✅ ABI manually extracted and exported${NC}"
    else
        echo -e "${RED}❌ Failed to extract ABI${NC}"
        exit 1
    fi
fi

# Validate the ABI JSON
echo -e "${YELLOW}🔍 Validating ABI JSON...${NC}"

if command -v jq &> /dev/null; then
    if jq empty "$CONTRACTS_DIR/abi.json" 2>/dev/null; then
        echo -e "${GREEN}✅ ABI JSON is valid${NC}"
        
        # Show ABI functions
        echo -e "${BLUE}📋 Available functions:${NC}"
        jq -r '.[].name' "$CONTRACTS_DIR/abi.json" | sed 's/^/   • /'
    else
        echo -e "${RED}❌ Invalid ABI JSON${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  jq not found, skipping JSON validation${NC}"
    echo -e "${YELLOW}💡 Install jq for JSON validation: brew install jq${NC}"
fi