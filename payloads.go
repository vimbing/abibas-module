package adidasv2

type LoginPayload struct {
	User            string `json:"user"`
	Password        string `json:"password"`
	PersistentLogin bool   `json:"persistentLogin"`
	Recaptcha       string `json:"recaptcha"`
}

type AddToCartPayload []AtcVariant

type AtcVariant struct {
	ProductID            string `json:"product_id"`
	Quantity             int    `json:"quantity"`
	ProductVariationSku  string `json:"product_variation_sku"`
	ProductID0           string `json:"productId"`
	Size                 string `json:"size"`
	DisplaySize          string `json:"displaySize"`
	SpecialLaunchProduct bool   `json:"specialLaunchProduct"`
}

type AddressPayload struct {
	Customer        Customer        `json:"customer"`
	ShippingAddress ShippingAddress `json:"shippingAddress"`
	BillingAddress  BillingAddress  `json:"billingAddress"`
	LegalAcceptance LegalAcceptance `json:"legalAcceptance"`
}
type Customer struct {
	Email             string `json:"email"`
	ReceiveSmsUpdates bool   `json:"receiveSmsUpdates"`
}
type ShippingAddress struct {
	EmailAddress string `json:"emailAddress"`
	Country      string `json:"country"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Address1     string `json:"address1"`
	Zipcode      string `json:"zipcode"`
	City         string `json:"city"`
	PhoneNumber  string `json:"phoneNumber"`
	Address2     string `json:"address2"`
}
type BillingAddress struct {
	EmailAddress string `json:"emailAddress"`
	Country      string `json:"country"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Address1     string `json:"address1"`
	Zipcode      string `json:"zipcode"`
	City         string `json:"city"`
	PhoneNumber  string `json:"phoneNumber"`
	Address2     string `json:"address2"`
	VatCode      string `json:"vatCode"`
}
type Acceptances struct {
	DocumentName string `json:"documentName"`
	Accepted     bool   `json:"accepted"`
}
type LegalAcceptance struct {
	ScenarioName     string        `json:"scenarioName"`
	AcceptanceLocale string        `json:"acceptanceLocale"`
	Acceptances      []Acceptances `json:"acceptances"`
}

type OrderPayload struct {
	BasketID        string `json:"basketId"`
	IsAsync         bool   `json:"isAsync"`
	PaymentMethodID string `json:"paymentMethodId"`
	Token           string `json:"token"`
}
