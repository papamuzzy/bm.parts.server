package apiProm

type Group struct {
	Id                   int64        `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Image                string       `json:"image"`
	ParentGroupId        int64        `json:"parent_group_id"`
	NameMultilang        MultiLangStr `json:"name_multilang"`
	DescriptionMultilang MultiLangStr `json:"description_multilang"`
}

type MultiLangStr struct {
	Ru string `json:"ru"`
	Uk string `json:"uk"`
}

type Product struct {
	Id                   int64       `json:"id"`
	ExternalId           string      `json:"external_id"`
	Name                 string      `json:"name"`
	Sku                  string      `json:"sku"`
	KeyWords             string      `json:"keywords"`
	Presence             string      `json:"presence"`
	Price                float64     `json:"price"`
	MinimumOrderQuantity int64       `json:"minimum_order_quantity"`
	Discount             Discount    `json:"discount"`
	Prices               interface{} `json:"prices"`
	Currency             string      `json:"currency"`
	Description          string      `json:"description"`
	Group                struct {
		Id            int64        `json:"id"`
		Name          string       `json:"name"`
		NameMultilang MultiLangStr `json:"name_multilang"`
	} `json:"group"`
	Category struct {
		Id      int64  `json:"id"`
		Caption string `json:"caption"`
	} `json:"category"`
	MainImage            string       `json:"main_image"`
	Images               []Image      `json:"images"`
	SellingType          string       `json:"selling_type"`
	Status               string       `json:"status"`
	QuantityInStock      int64        `json:"quantity_in_stock"`
	MeasureUnit          string       `json:"measure_unit"`
	IsVariation          bool         `json:"is_variation"`
	VariationBaseId      int64        `json:"variation_base_id"`
	VariationGroupId     int64        `json:"variation_group_id"`
	DateModified         string       `json:"date_modified"`
	InStock              bool         `json:"in_stock"`
	Regions              interface{}  `json:"regions"`
	NameMultilang        MultiLangStr `json:"name_multilang"`
	DescriptionMultilang MultiLangStr `json:"description_multilang"`
}

type Image struct {
	Id           int64  `json:"id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Url          string `json:"url"`
}

type Discount struct {
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	DateStart string  `json:"date_start"`
	DateEnd   string  `json:"date_end"`
}

type Order struct {
	Id               int64          `json:"id"`
	DateCreated      string         `json:"date_created"`
	ClientFirstName  string         `json:"client_first_name"`
	ClientSecondName string         `json:"client_second_name"`
	ClientLastName   string         `json:"client_last_name"`
	ClientId         int64          `json:"client_id"`
	ClientNotes      string         `json:"client_notes"`
	Products         []OrderProduct `json:"products"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	Price            string         `json:"price"`
	FullPrice        string         `json:"full_price"`
	DeliveryOption   struct {
		Id              int64  `json:"id"`
		Name            string `json:"name"`
		ShippingService string `json:"shipping_service"`
	} `json:"delivery_option"`
	DeliveryProviderData struct {
		Provider             string `json:"provider"`
		Type                 string `json:"type"`
		SenderWarehouseId    string `json:"sender_warehouse_id"`
		RecipientWarehouseId string `json:"recipient_warehouse_id"`
		DeclarationNumber    string `json:"declaration_number"`
		UnifiedStatus        string `json:"unified_status"`
	} `json:"delivery_provider_data"`
	DeliveryAddress string  `json:"delivery_address"`
	DeliveryCost    float64 `json:"delivery_cost"`
	PaymentOption   struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"payment_option"`
	PaymentData struct {
		Type           string `json:"type"`
		Status         string `json:"status"`
		StatusModified string `json:"status_modified"`
	} `json:"payment_data"`
	Status                    string `json:"status"`
	StatusName                string `json:"status_name"`
	Source                    string `json:"source"`
	HasOrderPromoFreeDelivery bool   `json:"has_order_promo_free_delivery"`
	CpaCommission             struct {
		Amount     string `json:"amount"`
		IsRefunded bool   `json:"is_refunded"`
	} `json:"cpa_commission"`
	Utm struct {
		Medium   string `json:"medium"`
		Source   string `json:"source"`
		Campaign string `json:"campaign"`
	} `json:"utm"`
	DontCallCustomerBack bool `json:"dont_call_customer_back"`
	PsPromotion          struct {
		Name       string   `json:"name"`
		Conditions []string `json:"conditions"`
	} `json:"ps_promotion"`
	Cancellation struct {
		Title     string `json:"title"`
		Initiator string `json:"initiator"`
	} `json:"cancellation"`
}

type OrderProduct struct {
	Id            int64        `json:"id"`
	ExternalId    string       `json:"external_id"`
	Image         string       `json:"image"`
	Quantity      float64      `json:"quantity"`
	Price         string       `json:"price"`
	Url           string       `json:"url"`
	Name          string       `json:"name"`
	NameMultilang MultiLangStr `json:"name_multilang"`
	TotalPrice    string       `json:"total_price"`
	MeasureUnite  string       `json:"measure_unite"`
	Sku           string       `json:"sku"`
	CpaCommission struct {
		Amount string `json:"amount"`
	} `json:"cpa_commission"`
}

type Orders struct {
	Orders []Order `json:"orders"`
}
