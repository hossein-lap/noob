package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const serverPort = 443
var token string
var coin_base string = "USDT"

// token handler {{{
func tokenHandler() string {
    token := os.Getenv("NOBTEX_API_TOKEN")
    if token == "" {
        file, err := os.ReadFile(".env")
        if err == nil {
            tostr := strings.TrimSuffix(string(file), "\n")
            tokens := strings.Split(tostr, "=")
            token = tokens[1]
        } else {
            log.Fatalln("There is no API token!!")
        }
    }
    return token
}
// }}}

// get user wallet method {{{
func getUserWallet(name string) string {
    var balance string
    type Wallet struct {
        ID      int    `json:"id"`
        Balance string `json:"balance"`
        Blocked string `json:"blocked"`
    }

    type Response struct {
        Status  string             `json:"status"`
        Wallets map[string]Wallet  `json:"wallets"`
    }

    token := tokenHandler()
    name = strings.ToLower(name) // Nobitex usually expects lowercase symbols

    // POST body
    data := strings.NewReader("currencies=" + name)

    req, err := http.NewRequest("POST", "https://api.nobitex.ir/v2/wallets", data)
    if err != nil {
        log.Fatalln("Error creating request:", err)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Token " + token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalln("Error making request:", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln("Error reading response body:", err)
    }

    var walletResp Response
    err = json.Unmarshal(body, &walletResp)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        log.Fatalln("Raw body:", string(body)) // Help debug bad responses
    }

    if len(walletResp.Wallets) == 0 {
        balance = "0.0"
    }

    for _, w := range walletResp.Wallets {
        balance = w.Balance
    }

    return balance
}
// }}}

// get price method {{{
func getPriceNobitex(name string) string {
    var token string = tokenHandler()
    type Response struct {
        Status         string      `json:"status"`
        LastUpdate     int64       `json:"lastUpdate"`
        LastTradePrice string      `json:"lastTradePrice"`
        Bids           [][]string  `json:"bids"`
        Asks           [][]string  `json:"asks"`
    }

    name = strings.ToUpper(name)
    url := "https://api.nobitex.ir/v3/orderbook/" + name
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalln("Error:", err)
    }
    resp.Header.Add("Authorization", "Token " + token)
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln("Error reading body:", err)
    }

    var price Response
    err = json.Unmarshal(body, &price)
    if err != nil {
        log.Fatalln("Error unmarshalling JSON:", err)
    }

    // fmt.Printf("Status: %s\n", price.Status)
    // fmt.Printf("Last Update: %d\n", price.LastUpdate)
    return price.LastTradePrice

    // log.Fatalln("Bids:")
    // for _, bid := range price.Bids {
    //     fmt.Printf("  Price: %s, Quantity: %s\n", bid[0], bid[1])
    // }
    // log.Fatalln("Asks:")
    // for _, ask := range price.Asks {
    //     fmt.Printf("  Price: %s, Quantity: %s\n", ask[0], ask[1])
    // }
    //
    // // Print the response status and body
    // log.Fatalln("Response Status:", resp.Status)
    // log.Fatalln("Response Body:", string(body))


}
// }}}

// user favorite markets {{{
func UserFavMarket() {
    var token = "Token " + tokenHandler()
    type Response struct {
        Status   string      `json:"status"`
        Fav      []string    `json:"favoriteMarkets"`
    }

    url := "https://api.nobitex.ir/users/markets/favorite"
    req, err := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    stat := strings.Split(resp.Status, " ")
    if stat[0] != "200" {
        log.Fatalf("Respond: %s: %s\n", stat[0], stat[1])
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading body:", err)
        return
    }

    var market Response
    err = json.Unmarshal(body, &market)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return
    }
    if market.Status == "ok" {
        fmt.Println("Market | Price | Wallet")
        fmt.Println("----------------------")
        for _, coin := range market.Fav {

            coin_from := strings.ReplaceAll(coin, coin_base, "")
            coin_to := strings.ReplaceAll(coin, coin_from, "")
            coin_balance := getUserWallet(coin_from)
            if coin_balance != "0.0" {
                fmt.Printf("%s/%s | ", coin_from, coin_to)
                fmt.Printf("%s | ", getPriceNobitex(coin))
                fmt.Printf("%s\n", coin_balance)
                fmt.Println("----------------------")
            }

            // if coin_from == "USDT" || coin_from == "BTC"  { 
            //     fmt.Printf("From: %s\n", coin_from)
            //     fmt.Printf("To: %s\n", coin_to)
            //     getPriceNobitex(coin)
            //     getUserWallet(coin_from)
            //     fmt.Println("----------------------")
            // }

            // os.Exit(2)
            // fmt.Printf("Coin: %s\n", coin)

        }
    }

}
// }}}

