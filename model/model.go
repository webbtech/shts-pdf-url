package model

// UrlRequest struct
type UrlRequest struct {
	EstimateNumber *int    `json:"number" validate:"required"`
	RequestType    *string `json:"requestType" validate:"required"`
}
