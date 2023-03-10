# Translate bot answers

For translation we will use [Yandex Translate API](https://cloud.yandex.com/en-ru/services/translate). Every answer from the Bored API will be translated `en->ru` before returning it to Telegram. Requests to the Translate API will be made on behalf of the special service account. To avoid high costs for the `Translate API`, the number of translations will be limited to 100 per hour.

Base instruction is here https://cloud.yandex.ru/docs/translate/api-ref/authentication and https://cloud.yandex.ru/docs/iam/operations/iam-token/create-for-sa#via-jwt.

## Prepare the service account
Create a serivce account `dronebot`.
```bash
yc iam service-account create \
  --name boredbot \
  --description 'Allows the boredbot application to access the Yandex API'
```

Assign the required role to this account.
```bash
folder_id=$(yc config get folder-id)
account_id=$(yc iam service-account get --name 'boredbot' --format json | jq .id -r)
yc resource-manager folder add-access-binding \
  --id "${folder_id}" \
  --role 'ai.translate.user' \
  --subject "serviceAccount:${account_id}"
```

Get `auth_key` for that account.
```bash
yc iam key create --service-account-name boredbot --output key.json
cat key.json | base64 -w0
rm key.json
```

## How it will work

Secret `auth_key` will be provided to the service. Before each request to the `Translate API`, the cached `IAM_token` will be checked. If it exists, use it in the request according to the [API](https://cloud.yandex.ru/docs/translate/api-ref/authentication).
```bash
curl -X POST \
  -H "Authorization: Bearer ${iam_token}" \
  -d '{"sourceLanguageCode":"en","targetLanguageCode":"ru","texts":["apple"]}' \
  "https://translate.api.cloud.yandex.net/translate/v2/translate"
```

If the `IAM_token` is missing or it has expired, a new one must be created. First, the application will create a [JWT token](https://cloud.yandex.ru/docs/iam/operations/iam-token/create-for-sa#jwt-create) with `service_account_id` and sign it with `auth_key`. Then exchange it via https://iam.api.cloud.yandex.net/iam/v1/tokens to get `IAM_token`.
```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d "{\"jwt\": \"${jwt}\"}" \
    https://iam.api.cloud.yandex.net/iam/v1/tokens
```

Application will cache `IAM_token` for ~30m < 1h.
