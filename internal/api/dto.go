package api

type packSizesRequest struct {
	PackSizes []int `json:"packSizes"`
}

type packSizesResponse struct {
	PackSizes []int `json:"packSizes"`
}

type calculateRequest struct {
	Items int `json:"items"`
}

type packEntry struct {
	Pack     int `json:"pack"`
	Quantity int `json:"quantity"`
}

type calculateResponse struct {
	Items      int         `json:"items"`
	PackSizes  []int       `json:"packSizes"`
	Solution   []packEntry `json:"solution"`
	TotalItems int         `json:"totalItems"`
	ExtraItems int         `json:"extraItems"`
}
