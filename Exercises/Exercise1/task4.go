package main

import (
    "fmt"
    "runtime"
)

var i int = 0

func countUp() {
    for k := 0; k <= 1000000; k++ {
        i += 1
    }
}

func countDown() {
    for k := 0; k <= 1000000; k++ {
        i -= 1
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    go countUp()
    go countDown()
    fmt.Println(i)
}
