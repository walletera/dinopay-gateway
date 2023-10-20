package testdata

var WithdrawalCreated = []byte(
    `
{
  "type": "WithdrawalCreated",
  "data": {
     "id": "0ae1733e-7538-4908-b90a-5721670cb093",
     "user_id": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
     "psp_id": "dinopay",
     "external_id": null,
     "amount": 100,
     "currency": "USD",
     "status": "pending",
     "beneficiary": {
       "id": "2f98dbe7-72ab-4167-9be5-ecd3608b55e4",
       "description": "Richard Roe DinoPay account",
       "account": {
        "holder": "Richard Roe",
        "number": 0,
        "routing_key": "1200079635"
       }
     }
  }
}
`)
