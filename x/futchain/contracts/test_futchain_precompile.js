const { Web3 } = require('web3');
const fs = require('fs');
const path = require('path');

// Configuration
const RPC_URL = 'http://localhost:8545';
const PRECOMPILE_ADDRESS = '0x0000000000000000000000000000000000000807';

// Load ABI from file
const abiPath = path.join(__dirname, 'abi.json');
const PRECOMPILE_ABI = JSON.parse(fs.readFileSync(abiPath, 'utf8'));

// Initialize Web3
const web3 = new Web3(RPC_URL);

// Create contract instance
const futchainContract = new web3.eth.Contract(PRECOMPILE_ABI, PRECOMPILE_ADDRESS);

// Test data - will be populated from GetUnfinishedMatches
let TEST_MATCH_ID = 1;
let TEST_LEAGUE_ID = 1;
let TEST_TEAM_ID = 1;
let AVAILABLE_MATCH_IDS = [];

async function testConnection() {
  console.log('🔗 Testing connection...');
  try {
    const blockNumber = await web3.eth.getBlockNumber();
    const chainId = await web3.eth.getChainId();
    console.log(`✅ Connected to chain ${chainId}, block ${blockNumber}`);
    return true;
  } catch (error) {
    console.error('❌ Connection failed:', error.message);
    return false;
  }
}

async function testGetMatch(matchId = TEST_MATCH_ID) {
  console.log(`\n🏈 Testing getMatch(${matchId})...`);
  try {
    const result = await futchainContract.methods.getMatch(matchId).call();
    console.log('✅ Match data received:');
    console.log('  📊 Match Details:');
    console.log(`    ID: ${result.id}`);
    console.log(`    League ID: ${result.leagueId}`);
    console.log(`    Name: ${result.name}`);
    console.log(`    Time: ${result.time}`);
    console.log('  🏠 Home Team:');
    console.log(`    ID: ${result.homeId}`);
    console.log(`    Name: ${result.homeName}`);
    console.log(`    Score: ${result.homeScore}`);
    console.log('  🚗 Away Team:');
    console.log(`    ID: ${result.awayId}`);
    console.log(`    Name: ${result.awayName}`);
    console.log(`    Score: ${result.awayScore}`);
    console.log('  📊 Status:');
    console.log(`    Started: ${result.started}`);
    console.log(`    Finished: ${result.finished}`);
    console.log(`    Cancelled: ${result.cancelled}`);
    
    // Update test data with league and team IDs from this match
    if (result.leagueId && parseInt(result.leagueId.toString()) > 0) {
      TEST_LEAGUE_ID = parseInt(result.leagueId.toString());
      console.log(`  🎯 Using league ID ${TEST_LEAGUE_ID} for league tests`);
    }
    if (result.homeId && parseInt(result.homeId.toString()) > 0) {
      TEST_TEAM_ID = parseInt(result.homeId.toString());
      console.log(`  🎯 Using team ID ${TEST_TEAM_ID} for team tests`);
    }
    
    return result;
  } catch (error) {
    console.error('❌ getMatch failed:', error.message);
    return null;
  }
}

async function testGetLeague(leagueId = TEST_LEAGUE_ID) {
  console.log(`\n🏆 Testing getLeague(${leagueId})...`);
  try {
    const result = await futchainContract.methods.getLeague(leagueId).call();
    console.log('✅ League data received:');
    console.log(`  ID: ${result.id}`);
    console.log(`  Name: ${result.name}`);
    console.log(`  Group Name: ${result.groupName}`);
    return result;
  } catch (error) {
    console.error('❌ getLeague failed:', error.message);
    return null;
  }
}

async function testGetTeam(teamId = TEST_TEAM_ID) {
  console.log(`\n👥 Testing getTeam(${teamId})...`);
  try {
    const result = await futchainContract.methods.getTeam(teamId).call();
    console.log('✅ Team data received:');
    console.log(`  ID: ${result.id}`);
    console.log(`  Name: ${result.name}`);
    return result;
  } catch (error) {
    console.error('❌ getTeam failed:', error.message);
    return null;
  }
}

async function testGetUnfinishedMatches() {
  console.log('\n⏰ Testing getUnfinishedMatches()...');
  try {
    const result = await futchainContract.methods.getUnfinishedMatches().call();
    console.log('✅ Unfinished matches received:');
    console.log(`  Count: ${result.length}`);
    if (result.length > 0) {
      console.log(`  IDs: [${result.slice(0, 10).join(', ')}${result.length > 10 ? '...' : ''}]`);
      
      // Update global test data with real match IDs
      AVAILABLE_MATCH_IDS = result.map(id => parseInt(id.toString()));
      if (AVAILABLE_MATCH_IDS.length > 0) {
        TEST_MATCH_ID = AVAILABLE_MATCH_IDS[0];
        console.log(`  🎯 Using match ID ${TEST_MATCH_ID} for subsequent tests`);
      }
    }
    return result;
  } catch (error) {
    console.error('❌ getUnfinishedMatches failed:', error.message);
    return null;
  }
}

async function testGasEstimation() {
  console.log('\n⛽ Testing gas estimation...');
  
  const tests = [
    { name: 'getMatch', method: () => futchainContract.methods.getMatch(TEST_MATCH_ID) },
    { name: 'getLeague', method: () => futchainContract.methods.getLeague(TEST_LEAGUE_ID) },
    { name: 'getTeam', method: () => futchainContract.methods.getTeam(TEST_TEAM_ID) },
    { name: 'getUnfinishedMatches', method: () => futchainContract.methods.getUnfinishedMatches() }
  ];

  for (const test of tests) {
    try {
      const gas = await test.method().estimateGas();
      console.log(`  ${test.name}: ${gas} gas`);
    } catch (error) {
      console.log(`  ${test.name}: estimation failed - ${error.message}`);
    }
  }
}

async function testRawCalls() {
  console.log('\n🔧 Testing raw calls...');
  
  // Test each function with raw calls
  const tests = [
    {
      name: 'getMatch',
      data: futchainContract.methods.getMatch(TEST_MATCH_ID).encodeABI()
    },
    {
      name: 'getLeague', 
      data: futchainContract.methods.getLeague(TEST_LEAGUE_ID).encodeABI()
    },
    {
      name: 'getTeam',
      data: futchainContract.methods.getTeam(TEST_TEAM_ID).encodeABI()
    },
    {
      name: 'getUnfinishedMatches',
      data: futchainContract.methods.getUnfinishedMatches().encodeABI()
    }
  ];

  for (const test of tests) {
    try {
      const result = await web3.eth.call({
        to: PRECOMPILE_ADDRESS,
        data: test.data
      });
      console.log(`  ${test.name}:`);
      console.log(`    Call data: ${test.data}`);
      console.log(`    Raw result: ${result.slice(0, 66)}...`);
      console.log(`    Result length: ${result.length} chars`);
    } catch (error) {
      console.log(`  ${test.name}: raw call failed - ${error.message}`);
    }
  }
}

async function performanceTest() {
  console.log('\n🚀 Performance testing...');
  
  const iterations = 10;
  const tests = [
    { name: 'getMatch', fn: () => testGetMatch(TEST_MATCH_ID) },
    { name: 'getLeague', fn: () => testGetLeague(TEST_LEAGUE_ID) },
    { name: 'getTeam', fn: () => testGetTeam(TEST_TEAM_ID) },
    { name: 'getUnfinishedMatches', fn: () => testGetUnfinishedMatches() }
  ];

  for (const test of tests) {
    const times = [];
    console.log(`  Testing ${test.name} (${iterations} iterations)...`);
    
    for (let i = 0; i < iterations; i++) {
      const start = Date.now();
      try {
        await test.fn();
        times.push(Date.now() - start);
      } catch (error) {
        console.log(`    Iteration ${i + 1} failed: ${error.message}`);
      }
    }
    
    if (times.length > 0) {
      const avg = times.reduce((a, b) => a + b, 0) / times.length;
      const min = Math.min(...times);
      const max = Math.max(...times);
      console.log(`    Average: ${avg.toFixed(2)}ms, Min: ${min}ms, Max: ${max}ms`);
    }
  }
}

async function interactiveTest() {
  console.log('\n🎮 Testing with real data from unfinished matches...');
  
  // Use real match IDs from getUnfinishedMatches
  const testIds = AVAILABLE_MATCH_IDS.length > 0 ? AVAILABLE_MATCH_IDS.slice(0, 5) : [1, 2, 3];
  
  console.log(`Testing with match IDs: [${testIds.join(', ')}]`);
  
  for (const matchId of testIds) {
    console.log(`\n--- Testing match ID ${matchId} ---`);
    const matchData = await testGetMatch(matchId);
    
    if (matchData) {
      // Test the league from this match
      const leagueId = parseInt(matchData.leagueId.toString());
      if (leagueId > 0) {
        await testGetLeague(leagueId);
      }
      
      // Test both teams from this match
      const homeTeamId = parseInt(matchData.homeId.toString());
      const awayTeamId = parseInt(matchData.awayId.toString());
      
      if (homeTeamId > 0) {
        await testGetTeam(homeTeamId);
      }
      if (awayTeamId > 0 && awayTeamId !== homeTeamId) {
        await testGetTeam(awayTeamId);
      }
    }
  }
}

async function debugInfo() {
  console.log('\n🔍 Debug Information:');
  console.log(`  Precompile Address: ${PRECOMPILE_ADDRESS}`);
  console.log(`  RPC URL: ${RPC_URL}`);
  console.log(`  ABI Functions: ${PRECOMPILE_ABI.map(f => f.name).join(', ')}`);
  
  // Check precompile address
  try {
    const code = await web3.eth.getCode(PRECOMPILE_ADDRESS);
    console.log(`  Code at address: ${code}`);
    
    const balance = await web3.eth.getBalance(PRECOMPILE_ADDRESS);
    console.log(`  Balance: ${web3.utils.fromWei(balance, 'ether')} ETH`);
  } catch (error) {
    console.log(`  Address check failed: ${error.message}`);
  }
}

async function runAllTests() {
  console.log('🔬 Futchain Precompile Test Suite');
  console.log('=====================================');
  
  // Test connection first
  const connected = await testConnection();
  if (!connected) {
    console.log('\n❌ Cannot proceed without connection');
    return;
  }

  // Run all tests - start with GetUnfinishedMatches to get real data
  await debugInfo();
  
  // First, get unfinished matches to populate test data
  console.log('\n🎯 Step 1: Getting real match data...');
  await testGetUnfinishedMatches();
  
  // Then test individual functions with real data
  console.log('\n🎯 Step 2: Testing with real match data...');
  await testGetMatch();
  await testGetLeague();
  await testGetTeam();
  
  // Performance and detailed testing
  console.log('\n🎯 Step 3: Performance and detailed testing...');
  await testGasEstimation();
  await testRawCalls();
  await performanceTest();
  await interactiveTest();
  
  console.log('\n🎉 All tests completed!');
}

// Export functions for individual testing
module.exports = {
  testConnection,
  testGetMatch,
  testGetLeague,
  testGetTeam,
  testGetUnfinishedMatches,
  testGasEstimation,
  testRawCalls,
  performanceTest,
  debugInfo,
  runAllTests,
  PRECOMPILE_ADDRESS,
  PRECOMPILE_ABI
};

// Run all tests if called directly
if (require.main === module) {
  runAllTests().catch(console.error);
}
