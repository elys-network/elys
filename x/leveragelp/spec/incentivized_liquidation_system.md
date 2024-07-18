
# Incentivized Liquidation System for Leverage Market

## Diagram

```plaintext
+------------+                  +------------+                  +-------------+                 
|            |                  |            |                  |             |                
| Depositors |                  | Borrowers  |                  | Liquidators |                
|            |                  |            |                  |             |                  
+-----+------+                  +-----+------+                  +-----+-------+               
      |                               |                               |                             
      | Deposit stable coins          | Borrow stable coins           | Check positions                
      | and earn interest             | to open leverage positions    | for liquidation or if certain price is met(stop loss)     
      |                               |                               |                              
+-----v------+                  +-----v------+                  +-----v-------+                   
|            |                  |            |                  |             |                 
|  Stable    |<-----------------| Leverage   |<-----------------| Liquidation |
|  Coin Pool |   Borrow coins   | Position   |   Monitor        |   Function  |   Monitor        
|  Module    |                  | Module     |   positions      |             |   positions       
+------------+                  +------------+                  +-------------+                  
      |                               |                               |                               
      |                               |                               |                                
      |  Interest                      | Position health calculated   | If health < threshold, or certain price is met(stop loss)     
      |                               | as lpAmount / debt            | trigger liquidation           
      |                               |                               | and receive reward             
      |                               |                               |                              
+-----v------+                  +-----v------+                  +-----v-------+                  
|            |                  |            |                  |             |                   
|  Depositor |                  |  Borrower  |                  |  Liquidator |                   
|  Earns     |                  | Manages    |                  | Earns       |                  
|  Interest  |                  | Position   |                  | Reward      |                  
+------------+                  +------------+                  +-------------+                   
```

## Explanation

- **Depositors** provide liquidity to the system and earn interest.
- **Borrowers** take leveraged positions by borrowing from the pool.
- **Liquidators** and **Bots** monitor positions and execute liquidations when necessary.
- The system rewards liquidators and bots for maintaining the health of the market.
- The **Liquidation Function** is responsible for managing the logic of liquidations and distributing rewards.

This approach distributes the load of monitoring and liquidation to incentivized users and bots, making the system more efficient and scalable.
