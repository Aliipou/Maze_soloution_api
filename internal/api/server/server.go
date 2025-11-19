package server

import (
	"context"
	"goapi/internal/api/handlers/data"
	"goapi/internal/api/handlers/device_config"
	"goapi/internal/api/handlers/maze_device"
	"goapi/internal/api/middleware"
	"goapi/internal/api/service"
	"log"
	"net/http"
)

type Server struct {
	ctx        context.Context
	HTTPServer *http.Server
	logger     *log.Logger
}

func NewServer(ctx context.Context, sf *service.ServiceFactory, logger *log.Logger) *Server {

	mux := http.NewServeMux()
	err := setupDataHandlers(mux, sf, logger)
	if err != nil {
		logger.Fatalf("Error setting up data handlers: %v", err)
	}

	err = setupMazeDeviceHandlers(mux, sf, logger)
	if err != nil {
		logger.Fatalf("Error setting up maze device handlers: %v", err)
	}

	err = setupDeviceConfigHandlers(mux, sf, logger)
	if err != nil {
		logger.Fatalf("Error setting up device config handlers: %v", err)
	}

	middlewares := []middleware.Middleware{
		middleware.BasicAuthenticationMiddleware,
		middleware.CommonMiddleware,
	}

	return &Server{
		ctx:    ctx,
		logger: logger,
		HTTPServer: &http.Server{
			Handler: middleware.ChainMiddleware(mux, middlewares...),
		},
	}
}

func (api *Server) Shutdown() error {
	api.logger.Println("Gracefully shutting down server...")
	return api.HTTPServer.Shutdown(api.ctx)
}

func (api *Server) ListenAndServe(addr string) error {
	api.HTTPServer.Addr = addr
	return api.HTTPServer.ListenAndServe()
}

// * REST API handlers for original data endpoint
func setupDataHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {

	ds, err := sf.CreateDataService(service.SQLiteDataService)
	if err != nil {
		return err
	}

	mux.HandleFunc("OPTIONS /*", func(w http.ResponseWriter, r *http.Request) {
		data.OptionsHandler(w, r)
	})
	mux.HandleFunc("POST /data", func(w http.ResponseWriter, r *http.Request) {
		data.PostHandler(w, r, logger, ds)
	})
	mux.HandleFunc("PUT /data", func(w http.ResponseWriter, r *http.Request) {
		data.PutHandler(w, r, logger, ds)
	})
	mux.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		data.GetHandler(w, r, logger, ds)
	})
	mux.HandleFunc("GET /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		data.GetByIDHandler(w, r, logger, ds)
	})
	mux.HandleFunc("DELETE /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		data.DeleteHandler(w, r, logger, ds)
	})
	return err
}

// * REST API handlers for maze device status
func setupMazeDeviceHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {

	mazeService, err := sf.CreateMazeDeviceStatusService(service.SQLiteDataService)
	if err != nil {
		return err
	}

	mux.HandleFunc("POST /device/status", func(w http.ResponseWriter, r *http.Request) {
		maze_device.PostHandler(w, r, logger, mazeService)
	})
	mux.HandleFunc("PUT /device/status", func(w http.ResponseWriter, r *http.Request) {
		maze_device.PutHandler(w, r, logger, mazeService)
	})
	mux.HandleFunc("GET /device/status", func(w http.ResponseWriter, r *http.Request) {
		maze_device.GetHandler(w, r, logger, mazeService)
	})
	mux.HandleFunc("GET /device/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		maze_device.GetByIDHandler(w, r, logger, mazeService)
	})
	mux.HandleFunc("DELETE /device/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		maze_device.DeleteHandler(w, r, logger, mazeService)
	})
	return nil
}

// * REST API handlers for device config
func setupDeviceConfigHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {

	configService, err := sf.CreateDeviceConfigService(service.SQLiteDataService)
	if err != nil {
		return err
	}

	mux.HandleFunc("POST /device/config", func(w http.ResponseWriter, r *http.Request) {
		device_config.PostHandler(w, r, logger, configService)
	})
	mux.HandleFunc("PUT /device/config", func(w http.ResponseWriter, r *http.Request) {
		device_config.PutHandler(w, r, logger, configService)
	})
	mux.HandleFunc("GET /device/config", func(w http.ResponseWriter, r *http.Request) {
		device_config.GetHandler(w, r, logger, configService)
	})
	mux.HandleFunc("GET /device/config/{id}", func(w http.ResponseWriter, r *http.Request) {
		device_config.GetByIDHandler(w, r, logger, configService)
	})
	mux.HandleFunc("DELETE /device/config/{id}", func(w http.ResponseWriter, r *http.Request) {
		device_config.DeleteHandler(w, r, logger, configService)
	})
	return nil
}
