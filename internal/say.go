package internal

import "fmt"

func Say(msg string) {
    fmt.Println(msg)
}

func Rotate(msg string) string {
    if len(msg) < 2 {
        return msg
    }
    return msg[1:] + string(msg[0])
}
