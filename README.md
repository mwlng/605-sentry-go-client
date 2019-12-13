# 605 sentry go client

## Usage
    package main

    import (
        "fmt"

        sentry "github.com/mwlng/605-sentry-go-client"
    )

    func main() {
        sentryCli := sentry.New()
        sentryCli.Login(
            "https://sentry.dev.605.nu",
            "<Your 605 sentry username>",
            "<Your 605 sentry password")

        fmt.Println(sentryCli.GetAccessToken(
            "dbc6de4289ffbeb0401e2e66f57f0e14a92b852d408a4e0c88fadd8e2310a93c", // Application client ID
            "token", // OAuth access token type
            "https://indxr.dev.605.nu"))
    }
