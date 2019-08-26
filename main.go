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
	} else {
		statusCode, err = tikvExec(ctx, cli, flag.Args())
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

		ts := tokenize(line)
		if len(ts) == 0 {
			continue
		}

		result := kvOps(ctx, cli, ts)
		if err := result.Err(); err != nil {
			fmt.Println(err.Error())
		} else {
			if _, ok := result.(*kvReadResult); ok {
				fmt.Println(result.Result())
			}
			if _, ok := result.(*kvWriteResult); ok {
				fmt.Println("successed!")
			}
		}
	}
}

func tikvExec(ctx context.Context, cli *rawkv.Client, args []string) (int, error) {
	ts := tokenizeFromArray(args)
	if len(ts) == 0 {
		return 0, nil
	}

	result := kvOps(ctx, cli, ts)
	if err := result.Err(); err != nil {
		return -1, err
	} else {
		if _, ok := result.(*kvReadResult); ok {
			fmt.Println(result.Result())
		}
		if _, ok := result.(*kvWriteResult); ok {
			fmt.Println("successed!")
		}
	}

	return 0, nil
}

func kvOps(ctx context.Context, cli *rawkv.Client, ts []Token) kvResult {
	var result kvResult

	switch ts[0].Type {
	case GET:
		result = kvGet(ctx, cli, ts[1:])
	case PUT:
		result = kvPut(ctx, cli, ts[1:])
	case DELETE:
		result = kvDelete(ctx, cli, ts[1:])
	default:
		result = &kvUnknownResult{Error: fmt.Errorf("No such command: %s\n", ts[0])}
	}

	return result
}

type kvResult interface {
	Result() string
	Err() error
}

type kvReadResult struct {
	ResultByte []byte
	Error      error
}

func (r *kvReadResult) Result() string {
	return string(r.ResultByte)
}

func (r *kvReadResult) Err() error {
	return r.Error
}

type kvWriteResult struct {
	Error error
}

func (r *kvWriteResult) Result() string {
	return ""
}

func (r *kvWriteResult) Err() error {
	return r.Error
}

type kvUnknownResult struct {
	Error error
}

func (r *kvUnknownResult) Result() string {
	return ""
}

func (r *kvUnknownResult) Err() error {
	return r.Error
}

func kvGet(ctx context.Context, cli *rawkv.Client, opds []Token) kvResult {
	if len(opds) != 1 {
		return &kvReadResult{ResultByte: nil, Error: fmt.Errorf("1 arg required")}
	}

	result, err := cli.Get(ctx, opds[0].Literal)

	return &kvReadResult{ResultByte: result, Error: err}
}

func kvPut(ctx context.Context, cli *rawkv.Client, opds []Token) kvResult {
	if len(opds) != 2 {
		return &kvWriteResult{Error: fmt.Errorf("2 arg required")}
	}

	return &kvWriteResult{Error: cli.Put(ctx, opds[0].Literal, opds[1].Literal)}
}

func kvDelete(ctx context.Context, cli *rawkv.Client, opds []Token) kvResult {
	if len(opds) != 1 {
		return &kvWriteResult{Error: fmt.Errorf("1 arg required")}
	}

	return &kvWriteResult{Error: cli.Delete(ctx, opds[0].Literal)}
}
