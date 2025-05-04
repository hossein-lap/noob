package main

import (
    // "log"
    "os"
    _ "noob/backend"
    // "noob/ui"
)

func main() {
    // back.GetPriceNobitex("btcusdt")
    back.WriteLocal("file.json")
    // ui.Init()
    // log.Printf("%s: %s\n", "tok", "test")
}
