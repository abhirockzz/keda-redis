package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		log.Fatal("REDIS_HOST env variable missing")
	}

	listName := os.Getenv("REDIS_LIST")
	if listName == "" {
		log.Fatal("REDIS_LIST env variable missing")
	}

	password := os.Getenv("REDIS_PASSWORD")

	fmt.Println("connecting to Redis host ", host)

	opt := &redis.Options{
		Addr: host,
	}

	if password != "" {
		opt.Password = password
	}
	client := redis.NewClient(opt)

	_, err := client.Ping().Result()

	if err != nil {
		fmt.Println("failed to connect to Redis", err)
		return
	}
	fmt.Println("successfully connected to Redis", host)
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("failed to close client conn ", err)
			return
		}
		fmt.Println("closed client connection")

	}()

	go func() {
		fmt.Println("waiting for items...")
		for {
			items, err := client.BRPop(0*time.Second, listName).Result()
			if err != nil {
				fmt.Println("unable to fetch item from list", err)
				continue
			}
			fmt.Println("got item from list -", items[1])
			time.Sleep(1 * time.Second)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	fmt.Println("press ctrl+c to exit...")
	<-exit
	fmt.Println("program exited")
}
