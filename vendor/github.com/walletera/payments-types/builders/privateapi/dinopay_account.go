package privateapi

import (
    "github.com/walletera/payments-types/privateapi"
)

type DinopayAccountBuilder struct {
    institutionName privateapi.OptString
    institutionId   privateapi.OptString
    currency        privateapi.Currency
    accountHolder   string
    accountNumber   string
}

func NewDinopayAccountBuilder() *DinopayAccountBuilder {
    return &DinopayAccountBuilder{
        currency: privateapi.CurrencyUSD,
    }
}

func (b *DinopayAccountBuilder) WithInstitutionName(name string) *DinopayAccountBuilder {
    b.institutionName = privateapi.OptString{
        Value: name,
        Set:   true,
    }
    return b
}

func (b *DinopayAccountBuilder) WithInstitutionId(id string) *DinopayAccountBuilder {
    b.institutionId = privateapi.OptString{
        Value: id,
        Set:   true,
    }
    return b
}

func (b *DinopayAccountBuilder) WithCurrency(currency privateapi.Currency) *DinopayAccountBuilder {
    b.currency = currency
    return b
}

func (b *DinopayAccountBuilder) WithAccountHolder(accountHolder string) *DinopayAccountBuilder {
    b.accountHolder = accountHolder
    return b
}

func (b *DinopayAccountBuilder) WithAccountNumber(accountNumber string) *DinopayAccountBuilder {
    b.accountNumber = accountNumber
    return b
}

func (b *DinopayAccountBuilder) Build() privateapi.Account {
    return privateapi.Account{
        InstitutionName: b.institutionName,
        InstitutionId:   b.institutionId,
        Currency:        b.currency,
        AccountDetails: privateapi.AccountAccountDetails{
            OneOf: privateapi.NewDinopayAccountDetailsAccountAccountDetailsSum(
                privateapi.DinopayAccountDetails{
                    AccountType:   "dinopay",
                    AccountHolder: b.accountHolder,
                    AccountNumber: b.accountNumber,
                },
            ),
        },
    }
}
