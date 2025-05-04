package backend

import (
	// "log"
	"log"
	"os"
    "time"
)

func WriteLocal(name string) {
    info, err := os.Stat(name)
    if err != nil { 
        log.Fatal(err) 
    } 
    mod := time.Since(info.ModTime())

    // mod.Minutes()
    // mod.String()
    // info.Name()
    // time.Now()

    if mod.Hours() >= 1.00 {
        log.Println("You can do it")
    }
  
    // // Gives the size of the file in bytes 
    // size := info.Size() 
    // log.Println("Size of the file:", size) 
  
}

