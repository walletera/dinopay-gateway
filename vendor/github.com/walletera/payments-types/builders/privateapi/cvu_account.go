package privateapi

import (
    "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

type CVUAccountBuilder struct {
    InstitutionName privateapi.OptString
    InstitutionId   privateapi.OptString
    cvu             string
    alias           string
}

func NewCVUAccountBuilder() *CVUAccountBuilder {
    return &CVUAccountBuilder{}
}

func (b *CVUAccountBuilder) WithInstitutionName(name string) *CVUAccountBuilder {
    b.InstitutionName = privateapi.NewOptString(name)
    return b
}

func (b *CVUAccountBuilder) WithInstitutionId(id string) *CVUAccountBuilder {
    b.InstitutionId = privateapi.NewOptString(id)
    return b
}

func (b *CVUAccountBuilder) WithCVU(cvu string) *CVUAccountBuilder {
    b.cvu = cvu
    return b
}

func (b *CVUAccountBuilder) WithAlias(alias string) *CVUAccountBuilder {
    b.alias = alias
    return b
}

func (b *CVUAccountBuilder) Build() (privateapi.Account, werrors.WError) {
    var accountDetails privateapi.AccountAccountDetails
    if len(b.cvu) > 0 {
        accountDetails = privateapi.AccountAccountDetails{
            OneOf: privateapi.NewCvuAccountDetailsAccountAccountDetailsSum(
                privateapi.CvuAccountDetails{
                    AccountType: "cvu",
                    RoutingInfo: privateapi.CvuAccountDetailsRoutingInfo{
                        OneOf: privateapi.NewCvuCvuRoutingInfoCvuAccountDetailsRoutingInfoSum(
                            privateapi.CvuCvuRoutingInfo{
                                CvuRoutingInfoType: "cvu",
                                Cvu:                b.cvu,
                            },
                        ),
                    },
                },
            ),
        }
    } else if len(b.alias) > 0 {
        accountDetails = privateapi.AccountAccountDetails{
            OneOf: privateapi.NewCvuAccountDetailsAccountAccountDetailsSum(
                privateapi.CvuAccountDetails{
                    AccountType: "cvu",
                    RoutingInfo: privateapi.CvuAccountDetailsRoutingInfo{
                        OneOf: privateapi.NewAliasCvuRoutingInfoCvuAccountDetailsRoutingInfoSum(
                            privateapi.AliasCvuRoutingInfo{
                                CvuRoutingInfoType: "alias",
                                Alias:              b.alias,
                            },
                        ),
                    },
                },
            ),
        }
    } else {
        return privateapi.Account{}, werrors.NewNonRetryableInternalError("either cvu or alias must be set in a cvu account")
    }
    return privateapi.Account{
        InstitutionName: b.InstitutionName,
        InstitutionId:   b.InstitutionId,
        Currency:        privateapi.CurrencyARS,
        AccountDetails:  accountDetails,
    }, nil
}
