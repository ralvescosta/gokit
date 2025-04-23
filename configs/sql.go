// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// SQLConfigs defines configuration parameters for SQL database connections.
// It contains all the necessary information to establish and maintain
// connections to SQL databases like PostgreSQL, MySQL, etc.
type SQLConfigs struct {
	// Host specifies the database server hostname or IP address
	Host string
	// Port defines the network port on which the database server is listening
	Port string
	// User specifies the username for database authentication
	User string
	// Password contains the authentication credential for the database user
	Password string
	// DbName specifies the name of the database to connect to
	DbName string
	// SecondsToPing defines the interval in seconds between health check pings
	// to verify the database connection remains active
	SecondsToPing int
}
