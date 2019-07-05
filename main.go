package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tikv/client-go/config"
	"github.com/tikv/client-go/rawkv"
)

const (
	StatusClientCreateFailed = 10
	StatusScanFailed         = 30
)

var (
	pdHosts string
)

func init() {
	flag.StringVarP(&pdHosts, "pdhosts", "h", "0.0.0.0:2379", "hosts of placement driver delimited by ','")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Default()

	cli, err := rawkv.NewClient(ctx, strings.Split(pdHosts, ","), cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(StatusClientCreateFailed)
	}

	var statusCode int
	// start interactive shell
	if len(flag.Args()) == 0 {
		statusCode, err = tikvShell(ctx, cli)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	os.Exit(statusCode)
}

func tikvShell(ctx context.Context, cli *rawkv.Client) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sc := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanned := sc.Scan()
		if !scanned {
			return StatusScanFailed, fmt.Errorf("failed to scan")
		}

		line := sc.Text()
		if line == "" {
			continue
		}

		ts := tokenizer(line)
		if len(ts) == 0 {
			continue
		}

		switch ts[0].Type {
		case GET:
			v, err := kvGet(ctx, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(string(v))
		case PUT:
			err := kvPut(ctx, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("successed!")
			}
		case DELETE:
			err := kvDelete(ctx, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("successed!")
			}
		default:
			fmt.Printf("No such command: %s\n", ts[0])
		}
	}
}

func kvGet(ctx context.Context, cli *rawkv.Client, opds []Token) ([]byte, error) {
	if len(opds) != 1 {
		return nil, fmt.Errorf("1 arg required")
	}

	return cli.Get(ctx, opds[0].Literal)
}

func kvPut(ctx context.Context, cli *rawkv.Client, opds []Token) error {
	if len(opds) != 2 {
		return fmt.Errorf("2 arg required")
	}

	return cli.Put(ctx, opds[0].Literal, opds[1].Literal)
}

func kvDelete(ctx context.Context, cli *rawkv.Client, opds []Token) error {
	if len(opds) != 1 {
		return fmt.Errorf("1 arg required")
	}

	return cli.Delete(ctx, opds[0].Literal)
}
