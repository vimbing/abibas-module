package adidasv2

import "time"

type LoginResponse struct {
	Acid                string `json:"acid"`
	UserName            string `json:"userName"`
	ExpirationTimestamp int64  `json:"expirationTimestamp"`
	IsPersistent        bool   `json:"isPersistent"`
}

type AddToCartResponse struct {
	BasketID     string    `json:"basketId"`
	Currency     string    `json:"currency"`
	ModifiedDate time.Time `json:"modifiedDate"`
	Pricing      struct {
		Total                            int     `json:"total"`
		BaseTotal                        int     `json:"baseTotal"`
		TotalTax                         float64 `json:"totalTax"`
		ProductTotal                     int     `json:"productTotal"`
		ProductTotalBeforeDiscounts      int     `json:"productTotalBeforeDiscounts"`
		ProductTotalBeforeOrderDiscounts int     `json:"productTotalBeforeOrderDiscounts"`
		ShippingTotal                    int     `json:"shippingTotal"`
		ShippingBaseTotal                int     `json:"shippingBaseTotal"`
	} `json:"pricing"`
	ResourceState     string `json:"resourceState"`
	TaxationPolicy    string `json:"taxationPolicy"`
	TotalProductCount int    `json:"totalProductCount"`
	MessageList       []struct {
		Type    string `json:"type"`
		Details struct {
			ShipmentID string `json:"shipmentId"`
		} `json:"details,omitempty"`
	} `json:"messageList"`
	ShipmentList []struct {
		ShipmentID          string `json:"shipmentId"`
		ShipmentType        string `json:"shipmentType"`
		ProductLineItemList []struct {
			ItemID               string `json:"itemId"`
			ProductID            string `json:"productId"`
			ProductName          string `json:"productName"`
			Category             string `json:"category"`
			CanonicalProductName string `json:"canonicalProductName"`
			ProductImage         string `json:"productImage"`
			Quantity             int    `json:"quantity"`
			Pricing              struct {
				BaseUnitPrice          int     `json:"baseUnitPrice"`
				UnitPrice              int     `json:"unitPrice"`
				BasePrice              int     `json:"basePrice"`
				Price                  int     `json:"price"`
				PriceAfterAllDiscounts int     `json:"priceAfterAllDiscounts"`
				UnitPriceWithoutTax    float64 `json:"unitPriceWithoutTax"`
			} `json:"pricing"`
			Gender         string `json:"gender"`
			Color          string `json:"color"`
			Size           string `json:"size"`
			AllowedActions struct {
				Delete         bool `json:"delete"`
				Edit           bool `json:"edit"`
				MoveToWishlist bool `json:"moveToWishlist"`
			} `json:"allowedActions"`
			MaxQuantityAllowed   int    `json:"maxQuantityAllowed"`
			IsBonusProduct       bool   `json:"isBonusProduct"`
			ProductType          string `json:"productType"`
			AvailableStock       int    `json:"availableStock"`
			LastAdded            bool   `json:"lastAdded"`
			IsFlashProduct       bool   `json:"isFlashProduct"`
			SpecialLaunchProduct bool   `json:"specialLaunchProduct"`
			ReturnType           int    `json:"returnType"`
		} `json:"productLineItemList"`
		ShippingLineItem struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			ID          string `json:"id"`
			Pricing     struct {
				BasePrice int `json:"basePrice"`
				Price     int `json:"price"`
			} `json:"pricing"`
			DiscountList []struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				CalloutMessage string `json:"calloutMessage"`
				CouponItemID   string `json:"couponItemId"`
			} `json:"discountList"`
			CarrierServiceName       string `json:"carrierServiceName"`
			FreeShippingThresholdMin int    `json:"freeShippingThresholdMin"`
			FreeShippingThresholdMax int    `json:"freeShippingThresholdMax"`
		} `json:"shippingLineItem"`
	} `json:"shipmentList"`
	Customer struct {
		CustomerID   string `json:"customerId"`
		CustomerEUCI string `json:"customerEUCI"`
		IsLoggedIn   bool   `json:"isLoggedIn"`
	} `json:"customer"`
	FreeShippingThreshold int `json:"freeShippingThreshold"`
}

type GetBasketsResponse struct {
	BasketID     string    `json:"basketId"`
	Currency     string    `json:"currency"`
	ModifiedDate time.Time `json:"modifiedDate"`
	Pricing      struct {
		Total                            int     `json:"total"`
		BaseTotal                        int     `json:"baseTotal"`
		TotalTax                         float64 `json:"totalTax"`
		ProductTotal                     int     `json:"productTotal"`
		ProductTotalBeforeDiscounts      int     `json:"productTotalBeforeDiscounts"`
		ProductTotalBeforeOrderDiscounts int     `json:"productTotalBeforeOrderDiscounts"`
		ShippingTotal                    int     `json:"shippingTotal"`
		ShippingBaseTotal                int     `json:"shippingBaseTotal"`
	} `json:"pricing"`
	ResourceState     string `json:"resourceState"`
	TaxationPolicy    string `json:"taxationPolicy"`
	TotalProductCount int    `json:"totalProductCount"`
	MessageList       []struct {
		Type string `json:"type"`
	} `json:"messageList"`
	ShipmentList []struct {
		ShipmentID          string `json:"shipmentId"`
		ShipmentType        string `json:"shipmentType"`
		ProductLineItemList []struct {
			ItemID               string `json:"itemId"`
			ProductID            string `json:"productId"`
			ProductName          string `json:"productName"`
			Category             string `json:"category"`
			CanonicalProductName string `json:"canonicalProductName"`
			ProductImage         string `json:"productImage"`
			Quantity             int    `json:"quantity"`
			Pricing              struct {
				BaseUnitPrice          int     `json:"baseUnitPrice"`
				UnitPrice              int     `json:"unitPrice"`
				BasePrice              int     `json:"basePrice"`
				Price                  int     `json:"price"`
				PriceAfterAllDiscounts int     `json:"priceAfterAllDiscounts"`
				UnitPriceWithoutTax    float64 `json:"unitPriceWithoutTax"`
			} `json:"pricing"`
			Gender         string `json:"gender"`
			Color          string `json:"color"`
			Size           string `json:"size"`
			AllowedActions struct {
				Delete         bool `json:"delete"`
				Edit           bool `json:"edit"`
				MoveToWishlist bool `json:"moveToWishlist"`
			} `json:"allowedActions"`
			MaxQuantityAllowed   int    `json:"maxQuantityAllowed"`
			IsBonusProduct       bool   `json:"isBonusProduct"`
			ProductType          string `json:"productType"`
			AvailableStock       int    `json:"availableStock"`
			IsFlashProduct       bool   `json:"isFlashProduct"`
			SpecialLaunchProduct bool   `json:"specialLaunchProduct"`
			ReturnType           int    `json:"returnType"`
		} `json:"productLineItemList"`
		ShippingLineItem struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			ID          string `json:"id"`
			Pricing     struct {
				BasePrice int `json:"basePrice"`
				Price     int `json:"price"`
			} `json:"pricing"`
			DiscountList []struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				CalloutMessage string `json:"calloutMessage"`
				CouponItemID   string `json:"couponItemId"`
			} `json:"discountList"`
			CarrierServiceName string `json:"carrierServiceName"`
			Delivery           struct {
				From time.Time `json:"from"`
				To   time.Time `json:"to"`
			} `json:"delivery"`
			FreeShippingThresholdMin int `json:"freeShippingThresholdMin"`
			FreeShippingThresholdMax int `json:"freeShippingThresholdMax"`
		} `json:"shippingLineItem"`
		ShippingOnDate string `json:"shippingOnDate"`
	} `json:"shipmentList"`
	ShippingAddress struct {
		Address1    string `json:"address1"`
		City        string `json:"city"`
		Country     string `json:"country"`
		FirstName   string `json:"firstName"`
		ID          string `json:"id"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
		Zipcode     string `json:"zipcode"`
	} `json:"shippingAddress"`
	BillingAddress struct {
		Address1      string `json:"address1"`
		City          string `json:"city"`
		Country       string `json:"country"`
		FirstName     string `json:"firstName"`
		ID            string `json:"id"`
		LastName      string `json:"lastName"`
		PhoneNumber   string `json:"phoneNumber"`
		Zipcode       string `json:"zipcode"`
		DocumentValue string `json:"documentValue"`
		VatCode       string `json:"vatCode"`
	} `json:"billingAddress"`
	Customer struct {
		CustomerID        string `json:"customerId"`
		Email             string `json:"email"`
		EncryptedEmail    string `json:"encryptedEmail"`
		CustomerEUCI      string `json:"customerEUCI"`
		ReceiveSmsUpdates bool   `json:"receiveSmsUpdates"`
		IsLoggedIn        bool   `json:"isLoggedIn"`
	} `json:"customer"`
	LegalAcceptance struct {
		ScenarioName     string `json:"scenarioName"`
		AcceptanceLocale string `json:"acceptanceLocale"`
		Acceptances      []struct {
			DocumentName string `json:"documentName"`
			Accepted     bool   `json:"accepted"`
		} `json:"acceptances"`
	} `json:"legalAcceptance"`
	FreeShippingThreshold int `json:"freeShippingThreshold"`
}

type OrderResponse struct {
	CheckoutID string `json:"checkoutId"`
	Parameters []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameters"`
	Basket struct {
		BasketID     string    `json:"basketId"`
		Currency     string    `json:"currency"`
		ModifiedDate time.Time `json:"modifiedDate"`
		Pricing      struct {
			Total                            int     `json:"total"`
			BaseTotal                        int     `json:"baseTotal"`
			TotalTax                         float64 `json:"totalTax"`
			ProductTotal                     int     `json:"productTotal"`
			ProductTotalBeforeDiscounts      int     `json:"productTotalBeforeDiscounts"`
			ProductTotalBeforeOrderDiscounts int     `json:"productTotalBeforeOrderDiscounts"`
			ShippingTotal                    int     `json:"shippingTotal"`
			ShippingBaseTotal                int     `json:"shippingBaseTotal"`
		} `json:"pricing"`
		ResourceState     string `json:"resourceState"`
		TaxationPolicy    string `json:"taxationPolicy"`
		TotalProductCount int    `json:"totalProductCount"`
		CheckoutID        string `json:"checkoutId"`
		ShipmentList      []struct {
			ShipmentID          string `json:"shipmentId"`
			ShipmentType        string `json:"shipmentType"`
			ProductLineItemList []struct {
				ItemID               string `json:"itemId"`
				ProductID            string `json:"productId"`
				ProductName          string `json:"productName"`
				CanonicalProductName string `json:"canonicalProductName"`
				ProductImage         string `json:"productImage"`
				Quantity             int    `json:"quantity"`
				Pricing              struct {
					BaseUnitPrice          int     `json:"baseUnitPrice"`
					UnitPrice              int     `json:"unitPrice"`
					BasePrice              int     `json:"basePrice"`
					Price                  int     `json:"price"`
					PriceAfterAllDiscounts int     `json:"priceAfterAllDiscounts"`
					UnitPriceWithoutTax    float64 `json:"unitPriceWithoutTax"`
				} `json:"pricing"`
				Gender         string `json:"gender"`
				Color          string `json:"color"`
				Size           string `json:"size"`
				AllowedActions struct {
				} `json:"allowedActions"`
				MaxQuantityAllowed   int    `json:"maxQuantityAllowed"`
				IsBonusProduct       bool   `json:"isBonusProduct"`
				ProductType          string `json:"productType"`
				AvailableStock       int    `json:"availableStock"`
				IsFlashProduct       bool   `json:"isFlashProduct"`
				SpecialLaunchProduct bool   `json:"specialLaunchProduct"`
				ReturnType           int    `json:"returnType"`
			} `json:"productLineItemList"`
			ShippingLineItem struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				ID          string `json:"id"`
				Pricing     struct {
					BasePrice int `json:"basePrice"`
					Price     int `json:"price"`
				} `json:"pricing"`
				DiscountList []struct {
					ID             string `json:"id"`
					Name           string `json:"name"`
					CalloutMessage string `json:"calloutMessage"`
					CouponItemID   string `json:"couponItemId"`
				} `json:"discountList"`
				CarrierServiceName string `json:"carrierServiceName"`
				Delivery           struct {
					From time.Time `json:"from"`
					To   time.Time `json:"to"`
				} `json:"delivery"`
				FreeShippingThresholdMin int `json:"freeShippingThresholdMin"`
				FreeShippingThresholdMax int `json:"freeShippingThresholdMax"`
			} `json:"shippingLineItem"`
			ShippingOnDate string `json:"shippingOnDate"`
		} `json:"shipmentList"`
		ShippingAddress struct {
			Address1    string `json:"address1"`
			City        string `json:"city"`
			Country     string `json:"country"`
			FirstName   string `json:"firstName"`
			ID          string `json:"id"`
			LastName    string `json:"lastName"`
			PhoneNumber string `json:"phoneNumber"`
			Zipcode     string `json:"zipcode"`
		} `json:"shippingAddress"`
		BillingAddress struct {
			Address1    string `json:"address1"`
			City        string `json:"city"`
			Country     string `json:"country"`
			FirstName   string `json:"firstName"`
			ID          string `json:"id"`
			LastName    string `json:"lastName"`
			PhoneNumber string `json:"phoneNumber"`
			Zipcode     string `json:"zipcode"`
		} `json:"billingAddress"`
		PaymentInstrumentList []struct {
			PaymentMethodID string `json:"paymentMethodId"`
			Amount          int    `json:"amount"`
			ID              string `json:"id"`
		} `json:"paymentInstrumentList"`
		Customer struct {
			CustomerID        string `json:"customerId"`
			Email             string `json:"email"`
			EncryptedEmail    string `json:"encryptedEmail"`
			CustomerEUCI      string `json:"customerEUCI"`
			ReceiveSmsUpdates bool   `json:"receiveSmsUpdates"`
			IsLoggedIn        bool   `json:"isLoggedIn"`
		} `json:"customer"`
		LegalAcceptance struct {
			ScenarioName     string `json:"scenarioName"`
			AcceptanceLocale string `json:"acceptanceLocale"`
			Acceptances      []struct {
				DocumentName string `json:"documentName"`
				Accepted     bool   `json:"accepted"`
			} `json:"acceptances"`
		} `json:"legalAcceptance"`
		FreeShippingThreshold int `json:"freeShippingThreshold"`
	} `json:"basket"`
}

type GetItemDataResponse struct {
	ID                 string `json:"id"`
	AvailabilityStatus string `json:"availability_status"`
	VariationList      []struct {
		Sku                string `json:"sku"`
		Size               string `json:"size"`
		Availability       int    `json:"availability"`
		AvailabilityStatus string `json:"availability_status"`
	} `json:"variation_list"`
}

type ScrapePaypalLinkResponse struct {
	Result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"result"`
	BuildNumber string `json:"buildNumber"`
	Timestamp   string `json:"timestamp"`
	Ndc         string `json:"ndc"`
	Redirect    struct {
		URL        string `json:"url"`
		Method     string `json:"method"`
		Parameters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters"`
		Target        string        `json:"target"`
		ShortURL      string        `json:"shortUrl"`
		ShopOrigin    string        `json:"shopOrigin"`
		Preconditions []interface{} `json:"preconditions"`
	} `json:"redirect"`
	AdditionalAttributes struct {
		ConnectorID string `json:"connectorId"`
	} `json:"additionalAttributes"`
	Workflow string `json:"workflow"`
}

type GetItemDataBackendResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ProductID       string `json:"product_id"`
	ModelNumber     string `json:"model_number"`
	OriginalPrice   int    `json:"original_price"`
	DisplayCurrency string `json:"display_currency"`
	Orderable       bool   `json:"orderable"`
	BadgeText       string `json:"badge_text"`
	BadgeColor      string `json:"badge_color"`
	Embedded        struct {
		Variations []struct {
			Size               string  `json:"size"`
			TechnicalSize      string  `json:"technical_size"`
			Orderable          bool    `json:"orderable"`
			AbsoluteSize       float64 `json:"absolute_size"`
			VariationProductID string  `json:"variation_product_id"`
			StockLevel         int     `json:"stock_level"`
			Links              struct {
				SimilarProducts struct {
					Href string `json:"href"`
				} `json:"similar_products"`
			} `json:"_links,omitempty"`
			LowOnStockMessage string `json:"low_on_stock_message,omitempty"`
		} `json:"variations"`
	} `json:"_embedded"`
	IsHype bool `json:"is_hype"`
}
