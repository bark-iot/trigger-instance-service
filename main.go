package main

import (
    "os"
    "log"
)

func main() {
    a := App{}
    a.Initialize(
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"))

    log.Println("Running on port 80")
    a.Run(":80")
}