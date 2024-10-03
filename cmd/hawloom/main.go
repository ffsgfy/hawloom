package main

import (
    "os"
    "os/signal"
    "time"

    "github.com/ffsgfy/hawloom/internal"
)

func main() {
    msg := "Hello, world! "

    tickerChan := time.Tick(time.Second)
    sigintChan := make(chan os.Signal, 1)
    signal.Notify(sigintChan, os.Interrupt)

    loop: for {
        select {
        case <-tickerChan:
            internal.Say(msg)
            msg = internal.Rotate(msg)
        case <-sigintChan:
            break loop
        }
    }
}
