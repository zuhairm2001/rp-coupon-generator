package models

type CouponRequest struct {
	Code             string `json:"code"`
	DiscountType     string `json:"discount_type"`
	Amount           int    `json:"amount"`
	IndividualUse    bool   `json:"individual_use"`
	ExcludeSaleItems bool   `json:"exclude_sale_items"`
	UsageLimit       int    `json:"usage_limit"`
	MinimumAmount    string `json:"minimum_amount"`
}

type CouponResponse struct {
	ID                       int        `json:"id"`
	Code                     string     `json:"code"`
	Amount                   string     `json:"amount"` // Changed from int to string
	Status                   string     `json:"status"` // Added missing field
	DateCreated              string     `json:"date_created"`
	DateCreatedGMT           string     `json:"date_created_gmt"`
	DateModified             string     `json:"date_modified"`
	DateModifiedGMT          string     `json:"date_modified_gmt"`
	DiscountType             string     `json:"discount_type"`
	Description              string     `json:"description"`
	DateExpires              *string    `json:"date_expires"`     // Changed to pointer since it can be null
	DateExpiresGMT           *string    `json:"date_expires_gmt"` // Changed to pointer since it can be null
	UsageCount               int        `json:"usage_count"`
	IndividualUse            bool       `json:"individual_use"`
	ProductIds               []int      `json:"product_ids"`
	ExcludedProductIds       []int      `json:"excluded_product_ids"`
	UsageLimit               *int       `json:"usage_limit"`            // Changed to pointer since it can be null
	UsageLimitPerUser        *int       `json:"usage_limit_per_user"`   // Changed to pointer since it can be null
	LimitUsageToXItems       *int       `json:"limit_usage_to_x_items"` // Changed to pointer since it can be null
	FreeShipping             bool       `json:"free_shipping"`
	ProductCategories        []int      `json:"product_categories"`
	ExcludeProductCategories []int      `json:"exclude_product_categories"`
	ExcludeSaleItems         bool       `json:"exclude_sale_items"`
	MinimumAmount            string     `json:"minimum_amount"`
	MaximumAmount            string     `json:"maximum_amount"`
	EmailRestrictions        []string   `json:"email_restrictions"`
	UsedBy                   []string   `json:"used_by"`
	MetaData                 []struct { // Added missing field
		ID    int    `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"meta_data"`
}
