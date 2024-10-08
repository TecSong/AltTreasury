# AltTreasury

AltTreasury is a blockchain-based fund management system for handling withdrawal requests, approvals, and executions.

## Table of Contents
1. [Quick Start](#quick-start)
2. [API Usage](#api-usage)
3. [Important Notes](#important-notes)

## Quick Start

### Requirements
- Docker and Docker Compose
- Go 1.22 or higher

### Installation and Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/TecSong/AltTreasury.git
   cd AltTreasury
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Start Geth and prepare the account:**
   1. Start Geth:
      ```bash
      docker-compose up -d geth
      ```
   2. Create a new account:
      1. Access the Geth container:
         ```bash
         docker exec -it geth /bin/sh
         ```
      2. Execute `geth account new` command using password `12345678` (or another password, but update it in docker-compose.yml)
      3. Copy the keystore file path (UTC--xxxx) to the `KEYSTORE_FILE_PATH` variable in docker-compose.yml
      4. Set the `TREASURY_ADDRESS` variable in docker-compose.yml to the newly created address

4. **Deploy the ERC20 token contract:**
   1. Use Remix IDE to connect to the local Geth node and deploy the ERC20 token contract
   2. Set the `TOKEN_ADDRESS` variable in docker-compose.yml to the deployed contract address

5. **Mint tokens:**
   - Mint some tokens to the account created in step 3
   - Transfer ETH to the account to cover gas fees

6. **Create database tables:**
   ```bash
   docker exec -it db /bin/sh
   mysql -u alt_treasury_admin -p alt_treasury < ddl.sql
   ```

7. **Start remaining services:**
   ```bash
   docker-compose up -d db api
   ```

## API Usage

1. **Create a withdrawal claim:**
   - Endpoint: `POST /v1/withdrawals/claim`
   - Request body example:
     ```json
     {
       "staff_id": 1001,
       "amount": 1,
       "recipient_address": "0x1234567890123456789012345678901234567890"
     }
     ```

2. **List withdrawal claims:**
   - Endpoint: `GET /v1/withdrawals/claims`
   - Query parameters: staff_id, status, created_after, created_before, page, page_size

3. **Get a specific withdrawal claim:**
   - Endpoint: `GET /v1/withdrawals/claims/{claim_id}`

4. **Approve a withdrawal claim:**
   - Endpoint: `POST /v1/withdrawals/claims/{claim_id}/approve`
   - Request body example:
     ```json
     {
       "manager_id": 1001
     }
     ```

5. **Reject a withdrawal claim:**
   - Endpoint: `POST /v1/withdrawals/claims/{claim_id}/reject`
   - Request body example:
     ```json
     {
       "manager_id": 1001
     }
     ```

6. **List withdrawal claim confirmations:**
   - Endpoint: `GET /v1/withdrawals/claims/confirmations`
   - Query parameters: staff_id, manager_id, action_type, page, page_size

## Important Notes

- The `ddl.sql` file is provided to create the necessary tables
- Manager IDs must be one of [1001, 1002, 1003] (simulating different roles)
- Ensure environment variables are correctly set and smart contracts are deployed before using the API
- All monetary operations use ERC20 tokens. Ensure accounts have sufficient token balance.
