package main

import(
    "fmt"
)

state string := backup
number int := 0
func main() {
    signal(idle)
}

}


switch signal {
    case backup:
        if signal == idle {
            // listen for heartbeat
            if noHeartbeat {
                signal(noHeartbeat)
        }
        if signal == noHeartbeat {
            state = primary
            signal(newPrimary)
        }
    case primary:
        if signal == newPrimary {
            number = readFile()
            signal(write)
        }
        if signal == write {
            number++
            writeFile(number)
            printNumber(number)
            signal(write)
        }

}
