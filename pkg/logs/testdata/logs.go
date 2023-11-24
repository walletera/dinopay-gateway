package testdata

const MultilineLog = `
2023/11/10 09:43:18 error processing message with payload: {
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
}: failed creating payment on dinopay: decode response: unexpected status code: 500
`

var Syslog = []string{
    "Nov 10 14:01:07 workstation kernel: [30755.510587] mce: CPU1: Package temperature/speed normal",
    "Nov 10 14:01:07 workstation kernel: [30755.510588] mce: CPU2: Package temperature/speed normal",
    "Nov 10 14:01:07 workstation kernel: [30755.510588] mce: CPU5: Package temperature/speed normal",
    "Nov 10 13:57:29 workstation wpa_supplicant[1494]: wlp60s0: CTRL-EVENT-BEACON-LOSS",
    "Nov 10 13:57:29 workstation gnome-shell[2780]: Window manager warning: Window 0x300e9d6 sets an MWM hint indicating it isn't resizable",
    "Nov 10 13:57:29 workstation gnome-shell[2780]: Window manager warning: but sets min size 1 x 1 and max size 2147483647 x 2147483647",
    "Nov 10 13:57:29 workstation gnome-shell[2780]: Window manager warning: this doesn't make much sense",
    "Nov 10 08:12:28 workstation /usr/lib/gdm3/gdm-x-session[2663]: (II) event1  - Power Button: device is a keyboard",
}

var SyslogSubstrs = []string{
    "CPU5: Package temperature/speed normal",
    "Window manager warning: this doesn't make much sense",
    "Power Button: device is a keyboard",
}
