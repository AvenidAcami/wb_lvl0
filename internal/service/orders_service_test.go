package service

import (
	"testing"
	"wb_lvl0/internal/model"
)

func TestOrderService_validateOrderInfo(t *testing.T) {
	type args struct {
		order model.Order
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				order: model.Order{
					OrderUid:    "b563feb7b2b84b6test",
					TrackNumber: "WBILMTESTTRACK",
					Entry:       "WBIL",
					Delivery: model.Delivery{
						Name:    "Test Testov",
						Phone:   "+9720000000",
						Zip:     "2639809",
						City:    "Kiryat Mozkin",
						Address: "Ploshad Mira 15",
						Region:  "Kraiot",
						Email:   "test@gmail.com",
					},
					Payment: model.Payment{
						Transaction:  "b563feb7b2b84b6test",
						RequestId:    "",
						Currency:     "USD",
						Provider:     "wbpay",
						Amount:       1817,
						PaymentDt:    1637907727,
						Bank:         "alpha",
						DeliveryCost: 1500,
						GoodsTotal:   317,
						CustomFee:    0,
					},
					Items: []model.Item{
						{
							ChrtId:      9934930,
							TrackNumber: "WBILMTESTTRACK",
							Price:       453,
							Rid:         "ab4219087a764ae0btest",
							Name:        "Mascaras",
							Sale:        30,
							Size:        "0",
							TotalPrice:  317,
							NmId:        2389212,
							Brand:       "Vivienne Sabo",
							Status:      202,
						},
					},
					Locale:            "en",
					InternalSignature: "",
					CustomerId:        "test",
					DeliveryService:   "meest",
					Shardkey:          "9",
					SmId:              99,
					DateCreated:       "2021-11-26T06:22:19Z",
					OofShard:          "1",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid (no items)",
			args: args{
				order: model.Order{
					OrderUid:    "b563feb7b2b84b6test",
					TrackNumber: "WBILMTESTTRACK",
					Entry:       "WBIL",
					Delivery: model.Delivery{
						Name:    "Test Testov",
						Phone:   "+9720000000",
						Zip:     "2639809",
						City:    "Kiryat Mozkin",
						Address: "Ploshad Mira 15",
						Region:  "Kraiot",
						Email:   "test@gmail.com",
					},
					Payment: model.Payment{
						Transaction:  "b563feb7b2b84b6test",
						RequestId:    "",
						Currency:     "USD",
						Provider:     "wbpay",
						Amount:       1817,
						PaymentDt:    1637907727,
						Bank:         "alpha",
						DeliveryCost: 1500,
						GoodsTotal:   317,
						CustomFee:    0,
					},
					Items:             []model.Item{},
					Locale:            "en",
					InternalSignature: "",
					CustomerId:        "test",
					DeliveryService:   "meest",
					Shardkey:          "9",
					SmId:              99,
					DateCreated:       "2021-11-26T06:22:19Z",
					OofShard:          "1",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid (some fields are empty)",
			args: args{
				order: model.Order{
					OrderUid:    "b563feb7b2b84b6test",
					TrackNumber: "",
					Entry:       "",
					Delivery: model.Delivery{
						Name:    "",
						Phone:   "+9720000000",
						Zip:     "2639809",
						City:    "Kiryat Mozkin",
						Address: "Ploshad Mira 15",
						Region:  "Kraiot",
						Email:   "test@gmail.com",
					},
					Payment: model.Payment{
						Transaction:  "b563feb7b2b84b6test",
						RequestId:    "",
						Currency:     "USD",
						Provider:     "wbpay",
						Amount:       1817,
						PaymentDt:    1637907727,
						Bank:         "",
						DeliveryCost: 1500,
						GoodsTotal:   317,
						CustomFee:    0,
					},

					Items: []model.Item{
						{
							ChrtId:      9934930,
							TrackNumber: "WBILMTESTTRACK",
							Price:       453,
							Rid:         "ab4219087a764ae0btest",
							Name:        "Mascaras",
							Sale:        30,
							Size:        "0",
							TotalPrice:  317,
							NmId:        2389212,
							Brand:       "Vivienne Sabo",
							Status:      202,
						},
					},
					Locale:            "en",
					InternalSignature: "",
					CustomerId:        "test",
					DeliveryService:   "meest",
					Shardkey:          "9",
					SmId:              99,
					DateCreated:       "2021-11-26T06:22:19Z",
					OofShard:          "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv := &OrderService{}
			if err := serv.validateOrderInfo(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("validateOrderInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
