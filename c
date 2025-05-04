#!/bin/bash

NOBITEX_API_TOKEN='bee506a52f28a7261d9a3dafe0a7ad00992a2ed9'
token=${NOBITEX_API_TOKEN}
alias curl='curl --silent'
curl='curl --silent'

userinfo() {
    $curl 'https://api.nobitex.ir/users/profile' \
      --header "Authorization: Token ${token}"
}

userwallet() {
    currency="$(echo $* | tr ' ' , | tr -d '\n')"
    $curl "https://api.nobitex.ir/v2/wallets?currencies=${currency}" \
      --header "Authorization: Token ${token}" | jq
}

userdeposite() {
    currency="{'currency': '${1}'}"
    $curl 'https://api.nobitex.ir/users/wallets/balance' \
      --request POST \
      --header "Authorization: Token ${token}" \
      --data "${currency}"
}

usertransaction() {
    $curl "https://api.nobitex.ir/users/wallets/transactions/list?wallet=${id}" \
      --header "Authorization: Token ${token}"
}

userfavmarket() {
    curl 'https://api.nobitex.ir/users/markets/favorite' \
      --header "Authorization: Token ${token}" \
      --request GET

}

# # usertransaction
# userfavmarket | jq | tee user-fav-market.json
userwallet usdt btc ltc | jq | tee user-wallet.json

# for id in $(jq -r '.favoriteMarkets[]' < user-fav-market.json); do
#     curl "https://api.nobitex.ir/v3/orderbook/${id}" \
#         | jq
# done

# curl https://api.nobitex.ir/v3/orderbook/DOGEUSDT \
#       --header "Authorization: Token ${token}"


