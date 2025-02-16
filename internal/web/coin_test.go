package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/service"
	svcmocks "github.com/miles0wu/meme-coin-api/internal/service/mocks"
	"github.com/miles0wu/meme-coin-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCoinHandler_Create(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.CoinService

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody Result
	}{
		{
			name: "create success",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "demo",
					Description: "desc",
				}).Return(domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
					PopularityScore: 0,
				}, nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"name": "demo", "description": "desc"}`))
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusCreated,
			wantBody: Result{
				Code: 200,
				Data: CoinVo{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       time.Now().Format(time.DateTime),
					UpdatedAt:       time.Now().Format(time.DateTime),
					PopularityScore: 0,
				},
			},
		},
		{
			name: "duplicate name error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "duplicate name",
					Description: "desc",
				}).Return(domain.Coin{}, service.ErrDuplicateName)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"name": "duplicate name", "description": "desc"}`))
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "coin name already exists",
			},
		},
		{
			name: "parse body error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"name": "duplicate name", "de`))
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid input",
			},
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "duplicate name",
					Description: "desc",
				}).Return(domain.Coin{}, errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"name": "duplicate name", "description": "desc"}`))
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinSvc := tc.mock(ctrl)
			// build handler
			hdl := NewCoinHandler(coinSvc, logger.NewNopLogger())

			// register route
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			bs, err := json.Marshal(tc.wantBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, string(bs), recorder.Body.String())
		})
	}
}

func TestCoinHandler_Detail(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.CoinService

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody Result
	}{
		{
			name: "get success",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
					PopularityScore: 0,
				}, nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodGet,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 200,
				Data: CoinVo{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       time.Now().Format(time.DateTime),
					UpdatedAt:       time.Now().Format(time.DateTime),
					PopularityScore: 0,
				},
			},
		},
		{
			name: "invalid id param",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodGet,
					"/api/v1/meme-coins/abc",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid id param",
			},
		},
		{
			name: "coin id not found",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{}, service.ErrNotFound)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodGet,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusNotFound,
			wantBody: Result{
				Code: 404,
				Msg:  "coin not found",
			},
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{}, errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodGet,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinSvc := tc.mock(ctrl)
			// build handler
			hdl := NewCoinHandler(coinSvc, logger.NewNopLogger())

			// register route
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			bs, err := json.Marshal(tc.wantBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, string(bs), recorder.Body.String())
		})
	}
}

func TestCoinHandler_Update(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.CoinService

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody Result
	}{
		{
			name: "update success",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				now := time.Now()
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}, nil)
				coinSvc.EXPECT().Update(gomock.Any(), domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc1",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}).Return(nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "desc1"}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/1",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 200,
				Msg:  "OK",
			},
		},
		{
			name: "parse body error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "de}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/1",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid input",
			},
		},
		{
			name: "invalid id param",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "desc1"}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/abc",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid id param",
			},
		},
		{
			name: "coin id not found",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{}, service.ErrNotFound)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "desc1"}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/1",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid id param",
			},
		},
		{
			name: "get db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{}, errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "desc1"}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/1",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
		{
			name: "update db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				now := time.Now()
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().GetById(gomock.Any(), int64(1)).Return(domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}, nil)
				coinSvc.EXPECT().Update(gomock.Any(), domain.Coin{
					Id:              1,
					Name:            "demo",
					Description:     "desc1",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}).Return(errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"description": "desc1"}`))
				req, err := http.NewRequest(
					http.MethodPut,
					"/api/v1/meme-coins/1",
					body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinSvc := tc.mock(ctrl)
			// build handler
			hdl := NewCoinHandler(coinSvc, logger.NewNopLogger())

			// register route
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			bs, err := json.Marshal(tc.wantBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, string(bs), recorder.Body.String())
		})
	}
}

func TestCoinHandler_Delete(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.CoinService

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody *Result
	}{
		{
			name: "delete success",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodDelete,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "invalid id param",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodDelete,
					"/api/v1/meme-coins/abc",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: &Result{
				Code: 400,
				Msg:  "invalid id param",
			},
		},
		{
			name: "coin id not found",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodDelete,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodDelete,
					"/api/v1/meme-coins/1",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinSvc := tc.mock(ctrl)
			// build handler
			hdl := NewCoinHandler(coinSvc, logger.NewNopLogger())

			// register route
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			// assertion
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantBody != nil {
				bs, err := json.Marshal(tc.wantBody)
				assert.NoError(t, err)
				assert.Equal(t, string(bs), recorder.Body.String())
			}
		})
	}
}

func TestCoinHandler_Poke(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.CoinService

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody Result
	}{
		{
			name: "poke success",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(nil)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins/1/poke",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 200,
				Msg:  "OK",
			},
		},
		{
			name: "invalid id param",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				return svcmocks.NewMockCoinService(ctrl)
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins/abc/poke",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Code: 400,
				Msg:  "invalid id param",
			},
		},
		{
			name: "coin id not found",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(service.ErrNotFound)
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins/1/poke",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Msg:  "invalid id param",
				Code: 400,
			},
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) service.CoinService {
				coinSvc := svcmocks.NewMockCoinService(ctrl)
				coinSvc.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/api/v1/meme-coins/1/poke",
					nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusInternalServerError,
			wantBody: Result{
				Code: 500,
				Msg:  "internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinSvc := tc.mock(ctrl)
			// build handler
			hdl := NewCoinHandler(coinSvc, logger.NewNopLogger())

			// register route
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			// assertion
			assert.Equal(t, tc.wantCode, recorder.Code)
			bs, err := json.Marshal(tc.wantBody)
			assert.NoError(t, err)
			assert.Equal(t, string(bs), recorder.Body.String())
		})
	}
}
