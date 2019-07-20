package pos

type invoiceMessage struct {
  RHash          string `json:"r_hash"`
  PaymentRequest string `json:"payment_request"`
  Settled        bool   `json:"settled"`
}

type invoiceStatusMessage struct {
  Settled bool `json:"settled"`
}
