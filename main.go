package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tikv/client-go/config"
	"github.com/tikv/client-go/rawkv"
)

func main() {
	ctx := context.Background()

	cfg := config.Default()
	cli, err := rawkv.NewClient(ctx, []string{"0.0.0.0:23791"}, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	statusCode, err := tikvShell(ctx, cli)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(statusCode)
}

func tikvShell(ctx context.Context, cli *rawkv.Client) (int, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	sc := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanned := sc.Scan()
		if !scanned {
			return 100, fmt.Errorf("failed to scan")
		}

		line := sc.Text()
		if line == "" {
			continue
		}

		ts := tokenizer(line)
		if len(ts) == 0 {
			continue
		}

		switch opr := string(ts[0]); opr {
		case "get":
			v, err := kvGet(ctx2, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(string(v))
		case "put":
			err := kvPut(ctx2, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("successed!")
			}
		case "delete":
			err := kvDelete(ctx2, cli, ts[1:])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("successed!")
			}
		default:
			fmt.Printf("No such command: %s\n", opr)
		}
	}
}

func kvGet(ctx context.Context, cli *rawkv.Client, opds []Token) ([]byte, error) {
	if len(opds) != 1 {
		return nil, fmt.Errorf("1 arg required")
	}

	return cli.Get(ctx, []byte(opds[0]))
}

func kvPut(ctx context.Context, cli *rawkv.Client, opds []Token) error {
	if len(opds) != 2 {
		return fmt.Errorf("2 arg required")
	}

	return cli.Put(ctx, []byte(opds[0]), []byte(opds[1]))
}

func kvDelete(ctx context.Context, cli *rawkv.Client, opds []Token) error {
	if len(opds) != 1 {
		return fmt.Errorf("1 arg required")
	}

	return cli.Delete(ctx, []byte(opds[0]))
}

//====

type Token []byte

func tokenizer(input string) []Token {
	ss := strings.Split(input, " ")

	tokens := make([]Token, len(ss))
	for i := range ss {
		lower := strings.ToLower(ss[i])
		tokens[i] = Token([]byte(lower))
	}

	return tokens
}
