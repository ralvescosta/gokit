// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// HTTPConfigs provides configuration settings for HTTP servers and clients.
// It contains parameters related to network binding, addressing, and diagnostic features.
type HTTPConfigs struct {
	// Host specifies the hostname or IP address to bind the HTTP server to
	Host string
	// Port defines the network port on which the HTTP server will listen
	Port string
	// Addr combines the Host and Port (typically in the format "host:port")
	// for direct use in network binding operations
	Addr string
	// EnableProfiling controls whether Go's pprof profiling endpoints should be exposed
	// for runtime debugging and performance analysis
	EnableProfiling bool
}
