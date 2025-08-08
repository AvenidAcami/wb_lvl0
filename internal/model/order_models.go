package model

type Payment struct {
	Transaction  string `json:"transaction" binding:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" binding:"required"`
	Provider     string `json:"provider" binding:"required"`
	Amount       int    `json:"amount" binding:"required"`
	PaymentDt    int    `json:"payment_dt" binding:"required"`
	Bank         string `json:"bank" binding:"required"`
	DeliveryCost int    `json:"delivery_cost" binding:"required"`
	GoodsTotal   int    `json:"goods_total" binding:"required"`
	CustomFee    int    `json:"custom_fee" binding:"required"`
}

type Delivery struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Zip     string `json:"zip" binding:"required"`
	City    string `json:"city" binding:"required"`
	Address string `json:"address" binding:"required"`
	Region  string `json:"region" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" binding:"required"`
	TrackNumber string `json:"track_number" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Rid         string `json:"rid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Sale        int    `json:"sale" binding:"required"`
	Size        string `json:"size" binding:"required"`
	TotalPrice  int    `json:"total_price" binding:"required"`
	NmId        int    `json:"nm_id" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	Status      int    `json:"status" binding:"required"`
}

type Order struct {
	OrderUid          string   `json:"order_uid" binding:"required"`
	TrackNumber       string   `json:"track_number" binding:"required"`
	Entry             string   `json:"entry" binding:"required"`
	Delivery          Delivery `json:"delivery" binding:"required"`
	Payment           Payment  `json:"payment" binding:"required"`
	Items             []Item   `json:"items" binding:"required"`
	Locale            string   `json:"locale" binding:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" binding:"required"`
	DeliveryService   string   `json:"delivery_service" binding:"required"`
	Shardkey          string   `json:"shardkey" binding:"required"`
	SmId              int      `json:"sm_id" binding:"required"`
	DateCreated       string   `json:"date_created" binding:"required"`
	OofShard          string   `json:"oof_shard" binding:"required"`
}
