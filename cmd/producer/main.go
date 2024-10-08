package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"k8s.io/klog/v2"

	"github.com/spf13/pflag"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

func init() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
}

// Connect to the database
func connectToDB() (*sql.DB, error) {
	dsn := "mysqluser:mysqlpw@tcp(mysql:3306)/inventory"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Check if the database connection is alive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	klog.Infof("connected to the database")
	return db, nil
}

// Insert a customer record into the database
func insertCustomer(db *sql.DB, firstName, lastName, email string) error {
	query := `INSERT INTO customers (first_name, last_name, email) VALUES (?, ?, ?)`
	_, err := db.Exec(query, firstName, lastName, email)
	if err != nil {
		return fmt.Errorf("failed to insert customer: %v", err)
	}
	klog.Infof("inserted: %s %s", firstName, lastName)
	return nil
}

func main() {
	klog.Infof("starting producer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to MySQL database
	db, err := connectToDB()
	if err != nil {
		klog.Fatalf("error connecting to the database: %v", err)
	}
	defer db.Close()

	wg := &sync.WaitGroup{}

	klog.Infof("run producer")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Start a goroutine to consume messages
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				firstName := "john"
				lastName := "doe"
				email := fmt.Sprintf(
					"john.doe.%d@example.com", time.Now().Unix())

				err := insertCustomer(
					db, firstName, lastName, email)
				if err != nil {
					klog.Errorf("error inserting customer: %v", err)
				}
			case <-ctx.Done():
				klog.Infof("exiting consumer, ctx done")
				return
			}
		}
	}()

	// graceful shutdown, no libs required, understand just below
	wg.Add(1)
	go func() {
		defer wg.Done()

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigterm:
			klog.Infof("sigterm received")
		case <-ctx.Done():
			klog.Infof("context done, bye")
			return
		}

		cancel()
	}()

	wg.Wait()
	klog.Infof("shutdown")
}
