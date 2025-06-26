package tradeshield

import (
	"database/sql"
	"fmt"
	"github.com/elys-network/elys/v6/indexer/db"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	OrderTypeLimitOpen  = int16(1)
	OrderTypeLimitClose = int16(2)
	LONG                = "LONG"
	SHORT               = "SHORT"
)

type PerpetualOrder struct {
	OwnerAddress     string    `json:"owner_address"`
	PoolID           int64     `json:"pool_id"`
	OrderID          int64     `json:"order_id"`
	OrderType        int16     `json:"order_type"`
	IsLong           bool      `json:"is_long"`
	Leverage         string    `json:"leverage"`
	CollateralAmount string    `json:"collateral_amount"`
	CollateralDenom  string    `json:"collateral_denom"`
	Price            string    `json:"price"`
	TakeProfitPrice  string    `json:"take_profit_price"`
	StopLossPrice    string    `json:"stop_loss_price"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreatePerpetualOrder inserts a new order into the database.
// This corresponds to the "CREATE" operation in CRUD.
func CreatePerpetualOrder(order *PerpetualOrder) error {
	// The SQL statement for insertion.
	// We use RETURNING to get the database-generated timestamps.
	sqlStatement := `
        INSERT INTO tradeshield.perpetual_orders (
            owner_address, pool_id, order_id, order_type, is_long, leverage, 
            collateral_amount, collateral_denom, price, take_profit_price, stop_loss_price
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING created_at, updated_at`

	// Execute the query
	err := db.DataBase.QueryRow(
		sqlStatement,
		order.OwnerAddress,
		order.PoolID,
		order.OrderID,
		order.OrderType,
		order.IsLong,
		order.Leverage,
		order.CollateralAmount,
		order.CollateralDenom,
		order.Price,
		order.TakeProfitPrice,
		order.StopLossPrice,
	).Scan(&order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		log.Printf("Error creating perpetual order: %v", err)
		return err
	}

	fmt.Printf("Successfully created order with ID %d for owner %s\n", order.OrderID, order.OwnerAddress)
	return nil
}

// GetPerpetualOrder retrieves an order from the database by its primary key.
// This corresponds to the "READ" operation in CRUD.
func GetPerpetualOrder(ownerAddress string, poolID int64, orderID int64) (*PerpetualOrder, error) {
	// The SQL statement for selection.
	sqlStatement := `
        SELECT owner_address, pool_id, order_id, order_type, is_long, leverage, 
               collateral_amount, collateral_denom, price, take_profit_price, 
               stop_loss_price, created_at, updated_at 
        FROM tradeshield.perpetual_orders 
        WHERE owner_address = $1 AND pool_id = $2 AND order_id = $3`

	order := &PerpetualOrder{}
	row := db.DataBase.QueryRow(sqlStatement, ownerAddress, poolID, orderID)

	// Scan the row data into the PerpetualOrder struct fields.
	err := row.Scan(
		&order.OwnerAddress,
		&order.PoolID,
		&order.OrderID,
		&order.OrderType,
		&order.IsLong,
		&order.Leverage,
		&order.CollateralAmount,
		&order.CollateralDenom,
		&order.Price,
		&order.TakeProfitPrice,
		&order.StopLossPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No order found with PK: %s, %d, %d", ownerAddress, poolID, orderID)
			return nil, fmt.Errorf("order not found")
		}
		log.Printf("Error getting perpetual order: %v", err)
		return nil, err
	}

	fmt.Printf("Successfully retrieved order with ID %d\n", order.OrderID)
	return order, nil
}

// UpdatePerpetualOrder updates an existing order in the database.
// This corresponds to the "UPDATE" operation in CRUD.
// This example updates price and stop loss.
func UpdatePerpetualOrder(ownerAddress string, poolID int64, orderID int64, newPrice, newStopLoss string) (*PerpetualOrder, error) {
	// The SQL statement for updating.
	// We also update the 'updated_at' timestamp to the current time.
	sqlStatement := `
        UPDATE tradeshield.perpetual_orders
        SET price = $4, stop_loss_price = $5, updated_at = CURRENT_TIMESTAMP
        WHERE owner_address = $1 AND pool_id = $2 AND order_id = $3
        RETURNING owner_address, pool_id, order_id, order_type, is_long, leverage,
                  collateral_amount, collateral_denom, price, take_profit_price, 
                  stop_loss_price, created_at, updated_at`

	updatedOrder := &PerpetualOrder{}

	// Execute the update and scan the returned row into the struct
	err := db.DataBase.QueryRow(sqlStatement, ownerAddress, poolID, orderID, newPrice, newStopLoss).Scan(
		&updatedOrder.OwnerAddress,
		&updatedOrder.PoolID,
		&updatedOrder.OrderID,
		&updatedOrder.OrderType,
		&updatedOrder.IsLong,
		&updatedOrder.Leverage,
		&updatedOrder.CollateralAmount,
		&updatedOrder.CollateralDenom,
		&updatedOrder.Price,
		&updatedOrder.TakeProfitPrice,
		&updatedOrder.StopLossPrice,
		&updatedOrder.CreatedAt,
		&updatedOrder.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No order found to update with PK: %s, %d, %d", ownerAddress, poolID, orderID)
			return nil, fmt.Errorf("order not found for update")
		}
		log.Printf("Error updating perpetual order: %v", err)
		return nil, err
	}

	fmt.Printf("Successfully updated order with ID %d\n", updatedOrder.OrderID)
	return updatedOrder, nil
}

// DeletePerpetualOrder removes an order from the database.
// This corresponds to the "DELETE" operation in CRUD.
func DeletePerpetualOrder(ownerAddress string, poolID int64, orderID int64) error {
	// The SQL statement for deletion.
	sqlStatement := `
        DELETE FROM tradeshield.perpetual_orders 
        WHERE owner_address = $1 AND pool_id = $2 AND order_id = $3`

	// Exec executes a query without returning any rows.
	res, err := db.DataBase.Exec(sqlStatement, ownerAddress, poolID, orderID)
	if err != nil {
		log.Printf("Error deleting perpetual order: %v", err)
		return err
	}

	// Check how many rows were affected.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no order found to delete")
	}

	fmt.Printf("Successfully deleted order with ID %d for owner %s\n", orderID, ownerAddress)
	return nil
}

// GetOrderBook retrieves a paginated list of orders based on the index, sorted by price.
func GetOrderBook(poolID int64, isLong bool, limit, offset int) ([]PerpetualOrder, error) {
	// The sorting order depends on whether we are fetching bids (longs) or asks (shorts).
	// Bids (longs) are typically sorted high to low (DESC).
	// Asks (shorts) are typically sorted low to high (ASC).
	sortOrder := "ASC" // For shorts/asks
	if isLong {
		sortOrder = "DESC" // For longs/bids
	}

	orderType := OrderTypeLimitOpen

	// Append LIMIT and OFFSET for pagination.
	// We add a secondary sort on created_at to ensure deterministic ordering (FIFO) for orders at the same price.
	sqlStatement := fmt.Sprintf(`
		SELECT owner_address, pool_id, order_id, order_type, is_long, leverage,
			   collateral_amount, collateral_denom, price, take_profit_price, 
			   stop_loss_price, created_at, updated_at 
		FROM tradeshield.perpetual_orders 
		WHERE pool_id = $1 AND order_type = $2 AND is_long = $3
		ORDER BY price %s, created_at ASC
		LIMIT $4 OFFSET $5`, sortOrder)

	rows, err := db.DataBase.Query(sqlStatement, poolID, orderType, isLong, limit, offset)
	if err != nil {
		log.Printf("Error querying order book: %v", err)
		return nil, err
	}
	defer rows.Close()

	var orders []PerpetualOrder
	for rows.Next() {
		var order PerpetualOrder
		err := rows.Scan(
			&order.OwnerAddress,
			&order.PoolID,
			&order.OrderID,
			&order.OrderType,
			&order.IsLong,
			&order.Leverage,
			&order.CollateralAmount,
			&order.CollateralDenom,
			&order.Price,
			&order.TakeProfitPrice,
			&order.StopLossPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning order book row: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating order book rows: %v", err)
		return nil, err
	}

	return orders, nil
}
