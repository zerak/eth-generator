package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	prefix    = flag.String("prefix", "", "The address prefix you expect [0-9][A-F][a-f]")
	suffix    = flag.String("suffix", "", "The address suffix you expect [0-9][A-F][a-f]")
	sensitive = flag.Bool("case", false, "Case sensitive default false")
	num       = flag.Int("num", 10, "Thread to use")

	ctx  = context.Background()
	quit = make(chan int)
)

func checkFlag() {
	flag.Parse()
	fmt.Printf("prefix:[%v] suffix:[%v] sensitive:%v num:%v\n", *prefix, *suffix, *sensitive, *num)

	check := func(str string) {
		if str == "" {
			return
		}
		match, err := regexp.MatchString(`^[0-9a-fA-F]{1,40}$`, str)
		if err != nil || !match {
			panic(fmt.Sprintf("check '%v' panic match:%v err:%v", str, match, err))
		}
	}
	check(*prefix)
	check(*suffix)
}

func process(ctx context.Context) {
	defer func() {
		quit <- 1
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("context down with error:%v\n", ctx.Err())
			return

		default:
			// Create an account
			key, err := crypto.GenerateKey()
			if err != nil {
				panic(err)
			}

			// Get the address
			address := crypto.PubkeyToAddress(key.PublicKey).Hex()
			privateKey := hex.EncodeToString(key.D.Bytes())
			fmt.Printf("address: %v key: %v\n", address, privateKey)

			var (
				pre = false
				suf = false
			)
			if *prefix != "" {
				preStr := "0x" + *prefix
				addrStr := address
				if !*sensitive {
					preStr = strings.ToUpper(preStr)
					addrStr = strings.ToUpper(addrStr)
				}
				pre = strings.HasPrefix(addrStr, preStr)
			}
			if *suffix != "" {
				sufStr := *suffix
				addrStr := address
				if !*sensitive {
					sufStr = strings.ToUpper(sufStr)
					addrStr = strings.ToUpper(addrStr)
				}
				suf = strings.HasSuffix(addrStr, sufStr)
			}

			if *prefix != "" && *suffix != "" {
				if pre && suf {
					fmt.Printf("with prefix:[%v] suffix:[%v]\n", *prefix, *suffix)
					fmt.Printf("find address: %v key: %v\n", address, privateKey)
					return
				} else {
					continue
				}
			}

			if *prefix != "" && pre {
				fmt.Printf("with prefix:[%v]\n", *prefix)
				fmt.Printf("find address: %v key: %v\n", address, privateKey)
				return
			}

			if *suffix != "" && suf {
				fmt.Printf("with suffix:[%v]\n", *suffix)
				fmt.Printf("find address: %v key: %v\n", address, privateKey)
				return
			}
		}
	}
}

func main() {
	fmt.Printf("ETH cool address generator.\n")

	checkFlag()

	// goroutine
	for i := 0; i < *num; i++ {
		go process(ctx)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	for {
		select {
		case <-quit:
			fmt.Printf("quit\n")
			return
		case sig := <-sigs:
			fmt.Printf("sig:%v\n", sig)
			return
		}
	}
}
