# **Vault Module Specification (V1)**

## **Overview**  
The `vault` module enables permissionless vault creation, where users deposit assets and receive shares. Vaults optimize returns using a manager that executes strategies without direct fund access. Performance tracking, fee distribution, and reinvestment mechanisms ensure sustainable growth.  

## **Key Features**  
- **Permissionless Vault Creation:** Anyone can deploy a vault.  
- **Deposits & Shares:** Users deposit assets and receive proportional vault shares.  
- **Redemptions:** Shares can be redeemed for increased value based on vault performance.  
- **Managed Actions:** A `manager` executes predefined actions securely.  
- **Performance Tracking:** Vaults record historical returns, yield, and risk metrics.  
- **Fee Structure:**  
  - **Performance Fee (X%)** deducted from profits.  
  - **Fee Distribution:**  
    - **Y% to Elys Treasury** for protocol sustainability.  
    - **Remaining to Manager** as an incentive for optimal returns.  
- **Reward Reinvestment:**  
  - Vaults **claim staking rewards** and reinvest them for compounding gains.  

---

## **Architecture**  
### **1. Core Components**  
1. **Vaults:**  
   - Store user deposits and issue shares.  
   - Maintain an internal balance sheet of assets and strategies.  
   - Track performance statistics.  

2. **Shares:**  
   - Represent ownership in the vault.  
   - Increase in value based on vault profits.  

3. **Manager:**  
   - Executes predefined actions (staking, swapping, LP provisioning, leveraging).  
   - Can be an AI agent or manually controlled.  
   - Earns a share of performance fees.  

4. **Performance Metrics:**  
   - **Total Value Locked (TVL)** – Tracks total funds in the vault.  
   - **Annual Percentage Yield (APY)** – Measures vault returns.  
   - **Historical Returns** – Records past performance for transparency.  

5. **Fee Mechanism:**  
   - **Performance Fee (X%)** deducted from net profits.  
   - **Fee Distribution:** Y% to Elys, remaining to the manager.  

6. **Reward Reinvestment:**  
   - **Claim Rewards:** Vault collects staking rewards.  
   - **Auto-Reinvest:** Rewards are automatically staked back for compounding.  

---

## **Flow Diagram**  
```
                 +------------------+
                 |   User Deposits   |
                 +------------------+
                          |
                          v
                 +------------------+
                 |  Vault Contract   |
                 +------------------+
                          |
                          v
                 +------------------+
                 |   Shares Issued   |
                 +------------------+
                          |
                          v
        +--------------------------------+
        |       Manager (AI/Manual)      |
        |  - Executes Actions            |
        |  - Optimizes Returns           |
        |  - Earns Performance Fees      |
        +--------------------------------+
                /       |       |        \
               v        v       v         v
  +----------------+ +----------------+ +----------------+
  |  Staking, Claim| |  Swapping       | |  Liquidity     |
  +----------------+ +----------------+ +----------------+
        |
        v
    +--------------------------------+
    | Vault Profits Accumulate       |
    | - Performance Tracked          |
    | - Fees Deducted (X%)           |
    +--------------------------------+
             |
    +------------------------+
    |    Fee Distribution    |
    | - Y% to Elys Treasury  |
    | - (1-Y)% to Manager    |
    +------------------------+
             |
             v
    +------------------------+
    |   Users Redeem Shares  |
    +------------------------+
```

---

## **Workflow**  
1. **Vault Creation:**  
   - Users create vaults permissionlessly.  
   - Vault parameters define strategy scope and fees.  

2. **Depositing Assets:**  
   - Users deposit assets and receive vault shares.  

3. **Investment Execution:**  
   - The manager executes staking, swaps, and LP provisioning.  
   - AI managers use historical data to optimize returns.  

4. **Reward Claim & Reinvestment:**  
   - Vault **claims rewards** from staking pools.  
   - Rewards **automatically reinvested** into staking.  
   - This **compounds returns** and increases APY.  

5. **Performance Tracking:**  
   - Vault records yield, APY, and profit/loss over time.  
   - Data available on-chain for transparency.  

6. **Fee Deduction:**  
   - Performance fees (X%) deducted from profits.  
   - Y% sent to Elys Treasury, remaining to the manager.  

7. **Redemptions:**  
   - Users burn shares to withdraw assets with accrued value.  


## **Security & Governance**  
- **Restricted Access:** Managers have no withdrawal rights, only allocation control.  
- **Auditable Performance:** Vault performance and fees are publicly verifiable.  
- **Governance Control:** Elys can update fee parameters via governance proposals.  


## **Future Scope**  
1. **Interchain Actions:**  
   - Expand vault actions using **Interchain Accounts (ICA)** to interact with external chains.  
   - Stake assets on remote chains for cross-chain yield farming.  
   - Swap assets on different blockchain networks.  

2. **Automated Risk Management:**  
   - Implement **AI-driven risk analysis** to prevent high-risk strategies.  
   - Dynamic allocation adjustments based on market conditions.  

3. **Dynamic Fee Structures:**  
   - Adaptive performance fee models based on manager performance.  
   - Lower fees for high-performing vaults to attract more users.  

4. **Insurance Mechanism:**  
   - Introduce a **vault insurance pool** to cover unexpected losses.  
   - Users can opt into insurance with a small fee for additional safety.