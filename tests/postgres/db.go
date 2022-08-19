package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tests/config"
	"log"
	"strings"
	"sync"
	"testing"
)

type TDB struct {
	Pool *pgxpool.Pool
	sync.Mutex
}

func NewFromEnv() *TDB {
	ctx := context.Background()
	config := config.GetConfig()

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DbHost(), config.DbPort(), config.DbUser(), config.DbPassword(), config.DbName())
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Ping database error: %v", err)
	}

	return &TDB{Pool: pool}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()

	d.Lock()
	d.Truncate(context.Background())
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate(context.Background())
}

func (d *TDB) Truncate(ctx context.Context) {
	var tables []string

	if err := pgxscan.Select(ctx, d.Pool, &tables, SelectAllTablesWithoutGoose); err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	if len(tables) == 0 {
		panic("run migration pls")
	}
	q := fmt.Sprintf("Truncate table %s RESTART IDENTITY", strings.Join(tables, ","))
	if _, err := d.Pool.Exec(ctx, q); err != nil {
		panic(err)
	}

}
