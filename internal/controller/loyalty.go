package controller

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/entity"
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
		router.Use(chiMW.AllowContentEncoding(GZip))
		// Собственная функция AllowContentType чтобы отдавать ошибку 400 Bad request
		router.Use(middleware.AllowContentType(ApplicationJSON))
		router.Use(middleware.Decompress)
		router.Use(compress)
		router.Post("/api/user/register", routes.PostRegUHandler)
		router.Post("/api/user/login", routes.PostLogUHandler)
	})

	// Private Routes
	// Only text/plain request type accepted
	// Only gzip request encoding accepted
	// If the client supports compression, the response will be compressed with gzip
	router.Group(func(router chi.Router) {
		router.Use(chiMW.AllowContentEncoding(GZip))
		// Собственная функция AllowContentType чтобы отдавать ошибку 400 Bad request
		router.Use(middleware.AllowContentType(TextPlain))
		router.Use(middleware.Decompress)
		router.Use(compress)
		router.Use(middleware.IdentifyUser)
		router.Post("/api/user/orders", routes.PostUOrder)
	})

	// router.Get("/api/user/orders", routes.GerUOrders)
	//	router.Get("/api/user/balance", routes.GetUBalance)
	//	router.Post("/api/user/balance/withdraw", routes.PostUWDBalance)
	//	router.Get("/api/user/withdrawals", routes.GetUWD)

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

// PostUOrder -.
func (routes *loyaltyRoutes) PostUOrder(w http.ResponseWriter, r *http.Request) {
	// var UserOrder *entity.UserOrder
	// Через контекст получаем reader
	// В случае необходимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	ctx := r.Context()
	reader := ctx.Value(middleware.KeyReader).(io.Reader)
	userID := ctx.Value(middleware.KeyUserID).(int64)

	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumberS := string(b)
	orderNumberI, err := strconv.Atoi(orderNumberS)
	if err != nil {
		http.Error(w, usecase.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
	}
	orderNumber := uint64(orderNumberI)

	_, err = routes.loyalty.PostUserOrder(ctx, &entity.UserOrder{
		UserID: userID,
		Number: orderNumber,
	})

	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOrderNumber) {
			http.Error(w, usecase.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, usecase.ErrOrderAlreadyRegByAnotherUser) {
			http.Error(w, usecase.ErrOrderAlreadyRegByAnotherUser.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, usecase.ErrOrderAlreadyRegByCurrUser) {
			w.Header().Set("Content-Type", ApplicationJSON)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusAccepted)
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
