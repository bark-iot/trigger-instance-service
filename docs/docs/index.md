# Trigger Instance Service

For full Bark documentation visit [http://localhost/docs](http://localhost/docs).

## Authorization

To perform any requests you need to send House Key and Signature via `Authorization` header. Example:
`Authorization: Bark <house_key>:<signature>`.
`signature` - is the sha256 hmac of the json body using House Secret as hash key.

## Create Action

POST `/trigger-instances`


*BODY JSON parameters*

Name          | Validation
------------  | -------------
device_token  | required
key           | required
input_data    | required

`key` - Trigger key. `input_data` - json string with key-value. `device_token` - Token of the Device.

*Response [200]*

```json
{
  "id": 1,
  "trigger_id": 1,
  "input_data": "[{\"key\":\"temp\",\"type\":\"int\"}]",
  "created_at": "2017-11-11 11:04:44 UTC"
}
```

*Error Response [400]*

Required fields are missing

*Error Response [403]*

Wrong authorization header

*Error Response [500]*

Server error