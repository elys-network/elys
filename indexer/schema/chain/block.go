package chain

import (
	"database/sql"
	"github.com/elys-network/elys/v6/indexer/db"
	"log"
	"time"
)

type Block struct {
	ID              string    `json:"id"`
	LastBlockHeight int64     `json:"last_block_height"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func UpsertBlockHeight(height int64) (*Block, error) {
	id := "latest_block_height"
	upsertBlockHeight := `
		INSERT INTO chain.block (id, last_block_height)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET last_block_height = EXCLUDED.last_block_height,
    		updated_at = CURRENT_TIMESTAMP
		RETURNING id, last_block_height, created_at, updated_at;
	`
	row := db.DataBase.QueryRow(upsertBlockHeight, id, height)
	block := &Block{}
	err := row.Scan(&block.ID, &block.LastBlockHeight, &block.CreatedAt, &block.UpdatedAt)
	if err != nil {
		return nil, err
	}
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
