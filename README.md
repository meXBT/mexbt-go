# meXBT go API client

Check http://godoc.org/github.com/meXBT/mexbt-go for function and response structs reference.

Check http://docs.mexbtpublicapi.apiary.io/ for public API documentation.

Check http://docs.mexbtprivateapisandbox.apiary.io/ for private API documentation.


# Private API example

    package main

    import "fmt"
    import "github.com/mexbt/mexbt-go"

    func main() {
        resp, err := mexbt.ProductPairs()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        if !resp.IsAccepted {
            fmt.Println("Server error:", resp.RejectReason)
            return
        }

        for _, pair := range resp.ProductPairs {
            fmt.Println("Got product pair named", pair.Name)
        }
    }

# Public API example

    package main

    import "fmt"
    import "github.com/mexbt/mexbt-go"

    var config = mexbt.Config{
                ApiKey:     "your_public_key",
                PrivateKey: "your_private_key",
                UserId:     "your@email",
    }

    func main() {

        // Uncomment to use sandbox API
        // mexbt.Sandbox = true
        resp, err := config.Balance()

        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        if !resp.IsAccepted {
            fmt.Println("Server error:", resp.RejectReason)
            return
        }

        for _, currency := range resp.Currencies {
            fmt.Printf("You have %v %v\n", currency.Balance, currency.Name)
        }
    }

