package controller

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/usecase"
)

const (
	ApplicationJSON = "application/json"
	TextPlain       = "text/plain"
	GZip            = "gzip"
)

type loyaltyRoutes struct {
	loyalty usecase.Loyalty
	cfg     *config.Config
}

// newLoyaltyRoutes -.
func newLoyaltyRoutes(router *chi.Mux, loyalty usecase.Loyalty, cfg *config.Config) *chi.Mux {
	routes := &loyaltyRoutes{
		loyalty: loyalty,
		cfg:     cfg,
	}

	// Если получилось, подключаем умную компрессию (только для ответа более 1400 байт)
	// иначе берем стандартное middleware chi для сжатия
	compress := defineCompression()

	// TODO Добавить middleware по обработке сжатых запросов

	// Public Routes
	// Only application/json request type accepted
	// Only gzip request encoding accepted
	// If the client supports compression, the response will be compressed with gzip
	router.Group(func(router chi.Router) {
		router.Use(middleware.AllowContentType(ApplicationJSON))
		router.Use(chiMW.AllowContentEncoding(GZip))
		router.Use(compress)
		router.Use(middleware.Decompress)
		router.Post("/api/user/register", routes.PostRegUHandler)
		router.Post("/api/user/login", routes.PostLogUHandler)
	})

	// Private Routes
	// Only application/json request type accepted
	// Only gzip request encoding accepted
	// If the client supports compression, the response will be compressed with gzip
	router.Group(func(router chi.Router) {
		router.Use(compress)
		router.Use(chiMW.AllowContentType(ApplicationJSON))
		router.Use(chiMW.AllowContentEncoding(GZip))
		//	router.Post("/api/user/orders", routes.PostUOrder)
		//	router.Get("/api/user/orders", routes.GerUOrders)
		//	router.Get("/api/user/balance", routes.GetUBalance)
		//	router.Post("/api/user/balance/withdraw", routes.PostUWDBalance)
		//	router.Get("/api/user/withdrawals", routes.GetUWD)

	})

	return router
}

// Если получилось, подключаем умную компрессию (только для ответа более 1400 байт)
// иначе берем стандартное middleware chi для сжатия
// defineCompression. -
func defineCompression() func(http.Handler) http.Handler {
	minSize := gziphandler.MinSize(1401)
	contentTypes := gziphandler.ContentTypes([]string{ApplicationJSON, TextPlain})
	compressionLevel := gziphandler.CompressionLevel(5)
	compress, err := gziphandler.GzipHandlerWithOpts(compressionLevel, minSize, contentTypes)
	if err != nil {
		compress = chiMW.Compress(5, ApplicationJSON, TextPlain)
	}

	return compress
}

// GerUOrders -.
func (routes *loyaltyRoutes) GerUOrders(w http.ResponseWriter, r *http.Request) {
	// UserOrders := entity.Loyalty{
	// 	UserID: "1",
	// 	UserOrders: []entity.UserOrder{
	// 		{Number: "1", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 		{Number: "2", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 		{Number: "3", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 	},
	// }
	//
	// userorders, _ := json.Marshal(UserOrders)
	//
	// w.Header().Set("Content-Type", ApplicationJSON)
	// w.WriteHeader(http.StatusOK)
	// w.Write(userorders)
}
