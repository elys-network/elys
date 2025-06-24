package chain

import (
	"database/sql"
	"fmt"
	"github.com/elys-network/elys/v6/indexer/db"
	"log"
	"time"
)

type Block struct {
	LastBlockHeight int64     `json:"last_block_height"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpsertBlock inserts a new block height record, or updates the timestamp if it already exists.
// This is more practical than a standard UPDATE since the only column is the PK.
func UpsertBlock(height int64) (*Block, error) {
	sqlStatement := `
		INSERT INTO chain.block (last_block_height)
		VALUES ($1)
		ON CONFLICT (last_block_height) DO UPDATE
		SET updated_at = CURRENT_TIMESTAMP
		RETURNING last_block_height, created_at, updated_at`

	block := &Block{}
	err := db.DataBase.QueryRow(sqlStatement, height).Scan(&block.LastBlockHeight, &block.CreatedAt, &block.UpdatedAt)
	if err != nil {
		log.Printf("Error upserting block record: %v", err)
		return nil, err
	}
	fmt.Printf("Successfully upserted block record for height %d\n", height)
	return block, nil
}

// GetLatestBlockHeight retrieves the highest block height from the table.
func GetLatestBlockHeight() (int64, error) {
	var height sql.NullInt64
	sqlStatement := `SELECT MAX(last_block_height) FROM chain.block`

	err := db.DataBase.QueryRow(sqlStatement).Scan(&height)
	if err != nil {
		log.Printf("Error getting latest block height: %v", err)
		return 0, err
	}

	if !height.Valid {
		return 0, nil
	}

	return height.Int64, nil
}
