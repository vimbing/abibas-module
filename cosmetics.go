package adidasv2

import "fmt"

func GetCosmetics(config *Config, atcResponse *AddToCartResponse) Cosmetics {
	emptyCosmetics := Cosmetics{
		Price: "UNKNOWN",
		Pid:   "UNKNOWN",
		Name:  "UNKNOWN",
		Image: "https://bitcoin.pl/wp-content/uploads/2021/11/kraken-shiba-inu-2.jpg",
		Size:  "UNKNOWN",
	}

	if len(atcResponse.ShipmentList) == 0 {
		return emptyCosmetics
	}

	if len(atcResponse.ShipmentList[0].ProductLineItemList) == 0 {
		return emptyCosmetics
	}

	product := atcResponse.ShipmentList[0].ProductLineItemList[0]

	return Cosmetics{
		Price: fmt.Sprint(product.Pricing.PriceAfterAllDiscounts),
		Pid:   config.DefaultConfig.TaskData.Sku,
		Name:  product.ProductName,
		Image: product.ProductImage,
		Size:  product.Size,
	}
}
