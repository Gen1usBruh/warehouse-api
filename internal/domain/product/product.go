package product

type Product struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
	Quantity    int32  `json:"quantity"`
}
