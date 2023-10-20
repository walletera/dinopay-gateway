package dinopay

// curl -X PUT 'localhost:1090/mockserver/openapi' \
//-d '{
//    "specUrlOrPayload": "https://raw.githubusercontent.com/walletera/dinopay/main/openapi/openapi.yaml",
//    "operationsAndResponses": {
//        "createPayment": "400"
//    }
//}'

var createPaymentFailed400Expectation = []byte(
    `
{
  "httpRequest" : {
    "operationId" : "createPayment",
    "specUrlOrPayload" : "https://raw.githubusercontent.com/walletera/dinopay/main/openapi/openapi.yaml"
  },
  "httpResponse" : {
    "statusCode" : 400,
    "headers" : {
      "content-type" : [ "text/html" ]
    },
    "body" : "some_string_value"
  },
  "id" : "00000000-13bf-c86c-0000-000013bfc86c",
  "priority" : 0,
  "timeToLive" : {
    "unlimited" : true
  },
  "times" : {
    "unlimited" : true
  }
}
`)
