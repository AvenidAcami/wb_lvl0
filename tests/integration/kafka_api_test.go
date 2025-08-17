package integration

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
	"wb_lvl0/config"
	"wb_lvl0/internal/controller"
	consumer "wb_lvl0/internal/kafka"
	"wb_lvl0/internal/repository"
	"wb_lvl0/internal/router"
	"wb_lvl0/internal/service"
)

func setupTestEnv(topic string) (*gin.Engine, *controller.OrdersController, *kafka.Writer) {
	r := gin.Default()

	// Инициализация компонентов
	config.InitENV()
	db := config.InitPostgres()
	redisPool := config.InitRedis()
	reader := config.InitKafka()

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, redisPool)
	orderController := controller.NewOrdersController(orderService)

	router.InitOrderRoutes(r, orderController)

	time.Sleep(5 * time.Second)
	// Запуск консьюмера
	go consumer.ParseOrders(reader, orderService)

	// Создание продюсера
	writer := kafka.Writer{
		Addr:  kafka.TCP("kafka:9092"),
		Topic: topic}

	return r, orderController, &writer
}

func sendKafkaMessage(t *testing.T, writer *kafka.Writer, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
	})
	if err != nil {
		t.Fatal("expected nil, got", err.Error())
	}
}

func sendAPIRequest(r *gin.Engine, t *testing.T, orderUID string) (string, map[string]interface{}) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/order/"+orderUID, nil)
	r.ServeHTTP(w, req)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("failed to unmarshal response:", err)
	}
	return strconv.Itoa(w.Code), resp
}

func TestKafkaAndApi(t *testing.T) {
	r, _, writer := setupTestEnv("orders")

	testCases := []struct {
		name      string
		orderJSON string
		orderUID  string
		wantCode  string
	}{
		{
			name:      "success",
			orderJSON: `{"order_uid":"53f81221-d160-459b-89e5-95bac37afd22","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"8d5365b2-84e0-47aa-a5ab-3fdd10af0152","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"8d5365b2-84e0-47aa-a5ab-3fdd10af0152","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`,
			orderUID:  "53f81221-d160-459b-89e5-95bac37afd22",
			wantCode:  "200",
		},
		{
			name:      "wrong uuid",
			orderJSON: "",
			orderUID:  "53",
			wantCode:  "400",
		},
		{
			name:      "not found",
			orderJSON: "",
			orderUID:  "53f81221-d160-459b-89e5-95bac37afd21",
			wantCode:  "404",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.orderJSON != "" {
				sendKafkaMessage(t, writer, tt.orderJSON)
				time.Sleep(5 * time.Second)
			}

			code, resp := sendAPIRequest(r, t, tt.orderUID)

			var expectedOrder map[string]interface{}
			if tt.orderJSON != "" {
				if err := json.Unmarshal([]byte(tt.orderJSON), &expectedOrder); err != nil {
					t.Fatal(err)
				}
			}

			if tt.wantCode != code {
				t.Fatal("expected code " + tt.wantCode + ", got " + code)
			} else {
				if (tt.wantCode == "200") && (!reflect.DeepEqual(resp["order"], expectedOrder)) {
					t.Fatalf("expected %+v, got %+v", expectedOrder, resp["order"])
				}
			}
		})
	}
}

// docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
// docker-compose -f docker-compose.test.yml down -v
