
# Incentivized Liquidation System for Leverage Market

## Diagram

```plaintext
+------------+                  +------------+                  +-------------+                 
|            |                  |            |                  |             |                
| Depositors |                  | Borrowers  |                  | Watchers    |                
|            |                  |            |                  |             |                  
+-----+------+                  +-----+------+                  +-----+-------+               
      |                               |                               |                             
      | Deposit stable coins          | Borrow stable coins           | Check positions                
      | and earn interest             | to open leverage positions    | for liquidation or if certain price is met(stop loss)     
      |                               |                               |                              
+-----v------+                  +-----v------+                  +-----v-------+                   
|            |                  |            |                  |             |                 
|  Stable    |<-----------------| Leverage   |<-----------------| Liquidation/|
|  Coin Pool |   Borrow coins   | Position   |   Monitor        |    Close    |   Monitor        
|  Module    |                  | Module     |   positions      |    Function |   positions       
+------------+                  +------------+                  +-------------+                  
      |                               |                               |                               
      |                               |                               |                                
      |  Interest                     | Position health calculated    | 1. If health < threshold, trigger liquidation
      |                               | as lpAmount / debt            | 2. certain price is met(stop loss), close position       
      |                               |                               | and receive reward             
      |                               |                               |                                  
+-----v------+                  +-----v------+                  +-----v-------+                  
|            |                  |            |                  |             |                   
|  Depositor |                  |  Borrower  |                  |  Watchers |                   
|  Earns     |                  | Manages    |                  | Earns       |                  
|  Interest  |                  | Position   |                  | Reward      |                  
+------------+                  +------------+                  +-------------+                   
```

## Explanation

- **Depositors** provide liquidity to the system and earn interest.
- **Borrowers** take leveraged positions by borrowing from the pool.
- **Watchers**: Watchers are users and bots that monitor positions for liquidation or stop loss conditions.
- **Bots**: Bots are automated programs specifically designed to efficiently handle these tasks.
- The system rewards liquidators and bots for maintaining the health of the market.
- The **Liquidation Function** is responsible for managing the logic of liquidations and distributing rewards.

This approach distributes the load of monitoring and liquidation to incentivized users and bots, making the system more efficient and scalable.
