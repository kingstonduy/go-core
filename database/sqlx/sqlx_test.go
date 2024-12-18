package sqlx

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/database"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func getConnection(t *testing.T) *database.Gdbc {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s", "postgres", "postgres", "localhost", 5432, "mydb", fmt.Sprintf("sslmode=%s", "disable"))

	gdbc, err := NewSqlxGdbc("postgres", dsn,
		database.WithMaxIdleCount(5),
		database.WithMaxOpen(5),
		database.WithMaxLifetime(60_000*time.Millisecond),
		database.WithMaxIdleTime(60_000*time.Millisecond),
		// database.WithLogger(logrus.NewLogrusLogger()),
	)

	if err != nil {
		t.Error(err)
	}

	return gdbc
}

func getConnectionWithLoggingHooks(t *testing.T) *database.Gdbc {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s", "postgres", "postgres", "localhost", 5432, "mydb", fmt.Sprintf("sslmode=%s", "disable"))

	gdbc, err := NewSqlxGdbc("postgres", dsn,
		database.WithMaxIdleCount(5),
		database.WithMaxOpen(5),
		database.WithMaxLifetime(60_000*time.Millisecond),
		database.WithMaxIdleTime(60_000*time.Millisecond),
		database.WithHooks(database.NewDBLogHooks()),
	)

	if err != nil {
		t.Error(err)
	}

	return gdbc
}

func TestGetConnection(t *testing.T) {
	getConnection(t)
}

func TestTransaction(t *testing.T) {
	gdbc := getConnectionWithLoggingHooks(t)
	ctx := context.Background()

	_, err := gdbc.Exec(ctx, "TRUNCATE TABLE product ")
	if err != nil {
		t.Error(err)
	}

	_ = gdbc.WithinTransaction(ctx, func(ctx context.Context) error {
		for i := 1; i <= 1; i++ {
			_, err := gdbc.Exec(ctx, "INSERT INTO product (id, name) VALUES ($1, $2)", i, "name")
			if err != nil {
				t.Error(err)
			}
			if i == 5 {
				return fmt.Errorf("Error from transaction ")
			}
		}
		return nil
	})

	time.Sleep(100 * time.Millisecond)

}

func TestConcurrency(t *testing.T) {
	gdbc := getConnection(t)

	ctx := context.Background()
	count := 100
	var (
		wg      sync.WaitGroup
		errChan = make(chan error)
	)

	wg.Add(count)

	for i := 1; i <= count; i++ {
		go func(i int, t *testing.T) {
			_, err := gdbc.Exec(ctx, "INSERT INTO product (id, name) VALUES ($1, $2)", i, "name")
			wg.Done()
			if err != nil {
				errChan <- err
			}

			status := gdbc.Stats(ctx)
			t.Log("STATUS:")
			t.Logf("Open connections: %v\n", status.OpenConnections)
			t.Logf("In Used connections: %v\n", status.InUse)
			t.Logf("Wait connections: %v\n", status.WaitCount)
			t.Logf("=======================================\n")
		}(i, t)
	}

	wg.Wait()
}

func TestPoolConnection(t *testing.T) {
	gdbc := getConnection(t)

	ctx := context.Background()

	_, err := gdbc.Exec(ctx, "INSERT INTO product (id, name) VALUES ($1, $2)", 10, "name")
	if err != nil {
		t.Logf("ERROR: %v", err)
	}

	status := gdbc.Stats(ctx)
	t.Log("STATUS:")
	t.Logf("Open connections: %v\n", status.OpenConnections)
	t.Logf("Free connections: %v\n", status.Idle)
	t.Logf("In Used connections: %v\n", status.InUse)
	t.Logf("Wait connections: %v\n", status.WaitCount)
	t.Logf("=======================================\n")

}

func TestTruncateTable(t *testing.T) {
	gdbc := getConnection(t)

	_, err := gdbc.Exec(context.Background(), "TRUNCATE TABLE product")

	if err != nil {
		t.Error(err)
	}

	var count int64
	rows, err := gdbc.Query(context.Background(), "SELECT COUNT(*) FROM product")
	if err != nil {
		t.Fatal(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			t.Fatal(err.Error())
		}
	}

	assert.Equal(t, int64(0), count)
}
