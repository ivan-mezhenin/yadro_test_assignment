package main

import (
	"flag"
	"fmt"
	"os"

	"dns-manager/client/config"
	grpcClient "dns-manager/client/internal/adapter/grpc"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file path")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		runAdd(cfg.GRPC.Address)
	case "remove":
		runRemove(cfg.GRPC.Address)
	case "list":
		runList(cfg.GRPC.Address)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runAdd(serverAddr string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnsctl add <ip> [--server addr]\n")
		os.Exit(1)
	}

	client, err := grpcClient.NewClient(*server)
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

func runRemove(serverAddr string) {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnsctl remove <ip> [--server addr]\n")
		os.Exit(1)
	}

	client, err := grpcClient.NewClient(*server)
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

func runList(serverAddr string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	server := fs.String("server", serverAddr, "gRPC server address")
	fs.Parse(os.Args[2:])

	client, err := grpcClient.NewClient(*server)
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
  --config path  Config file path (default: config.yaml)
  --server addr  gRPC server address (default: from config)`)
}
