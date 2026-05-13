package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"dns-manager/client/config"
	grpcClient "dns-manager/client/internal/adapter/grpc"
)

func main() {
	cfg := config.MustLoad("client/config/config.yaml")

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "add":
		runAdd(cfg.GRPC.Address, time.Duration(cfg.GRPC.Timeout)*time.Second)
	case "remove":
		runRemove(cfg.GRPC.Address, time.Duration(cfg.GRPC.Timeout)*time.Second)
	case "list":
		runList(cfg.GRPC.Address, time.Duration(cfg.GRPC.Timeout)*time.Second)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runAdd(serverAddr string, timeout time.Duration) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnsctl add <ip> [--server addr]\n")
		os.Exit(1)
	}

	client, err := grpcClient.NewClient(*server, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ok, errMsg := client.Add(args[0])
	if ok {
		fmt.Printf("added: %s\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", errMsg)
		os.Exit(1)
	}
}

func runRemove(serverAddr string, timeout time.Duration) {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnsctl remove <ip> [--server addr]\n")
		os.Exit(1)
	}

	client, err := grpcClient.NewClient(*server, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ok, errMsg := client.Remove(args[0])
	if ok {
		fmt.Printf("removed: %s\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", errMsg)
		os.Exit(1)
	}
}

func runList(serverAddr string, timeout time.Duration) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	client, err := grpcClient.NewClient(*server, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ok, errMsg, servers := client.List()
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", errMsg)
		os.Exit(1)
	}
	if len(servers) == 0 {
		fmt.Println("no dns servers configured")
		return
	}
	for _, s := range servers {
		fmt.Println(s)
	}
}

func printUsage() {
	fmt.Println(`CLI client for managing DNS servers

Commands:
  add <ip>       Add a DNS server
  remove <ip>    Remove a DNS server
  list           List all DNS servers
  help           Show this help

Options:
  --server addr  gRPC server address (default: localhost:50051)

Examples:
  dnsctl add 8.8.8.8
  dnsctl remove 8.8.8.8
   dnsctl list`)
}
