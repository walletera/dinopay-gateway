package dinopay

// curl -X PUT 'localhost:1090/mockserver/openapi' \
//-d '{
//    "specUrlOrPayload": "https://raw.githubusercontent.com/walletera/dinopay/main/openapi/openapi.yaml",
//    "operationsAndResponses": {
//        "createPayment": "201"
//    }
//}'

var CreatePaymentSucceedExpectation = []byte(
    `
{
  "id": "createPaymentSucceed",
  "httpRequest" : {
    "operationId" : "createPayment",
    "specUrlOrPayload" : "https://raw.githubusercontent.com/walletera/dinopay/main/openapi/openapi.yaml"
  },
  "httpResponse" : {
    "statusCode" : 201,
    "headers" : {
      "content-type" : [ "application/json" ]
    },
    "body" : {
      "id" : "bb17667e-daac-41f6-ada3-2c22f24caf22",
      "amount" : 100,
      "currency" : "USD",
      "sourceAccount" : {
        "accountHolder" : "john doe",
        "accountNumber" : "IE12BOFI90000112345678"
      },
      "destinationAccount" : {
        "accountHolder" : "jane doe",
        "accountNumber" : "IE12BOFI90000112349876"
      },
      "status" : "pending",
      "customerTransactionId" : "9713ec22-cf8d-4a21-affb-719db00d7388",
      "createdAt" : "2023-07-07",
      "updatedAt" : "2023-07-07"
    }
  },
  "priority" : 0,
  "timeToLive" : {
    "unlimited" : true
  },
  "times" : {
    "unlimited" : true
  }
}
`)
