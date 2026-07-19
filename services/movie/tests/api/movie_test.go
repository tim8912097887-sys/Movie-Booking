package movie_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/configs"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/db"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/interface/api"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/response"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/shutdown"
)

func getCreateMovieSchema(modifier ...func(*api.CreateMovieRequest)) *api.CreateMovieRequest {
	req := &api.CreateMovieRequest{
		Title:           "title",
		Description:     "description",
		Genres:          []string{"ACTION", "COMEDY"},
		DurationMinutes: 120,
		Rating:          "PG13",
		ReleaseDate:     time.Now(),
		DueDate:         time.Now().Add(time.Hour * 24),
	}
	for _, m := range modifier {
		if m == nil {
			continue
		}
		m(req)
	}
	return req
}

func getUpdateMovieSchema(modifier ...func(*api.UpdateMovieRequest)) *api.UpdateMovieRequest {
	req := &api.UpdateMovieRequest{
		Title:           nil,
		Description:     nil,
		Genres:          nil,
		DurationMinutes: nil,
		Rating:          nil,
		ReleaseDate:     nil,
		DueDate:         nil,
	}
	for _, m := range modifier {
		if m == nil {
			continue
		}
		m(req)
	}
	return req
}

func setupRouter(t *testing.T,h *api.Handler) *chi.Mux {
	t.Helper()
	r := chi.NewRouter()
	r.Route("/api/v1/movies",func(r chi.Router) {
		h.RegisterRoutes(r)
	})
	return r
}


func decodeResponse[T any](t *testing.T,resp *http.Response) T {
	t.Helper()
	var payload T
	err := json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}


func wireupHandler(t *testing.T) (*api.Handler,*pgxpool.Pool) {
	t.Helper()
	// Initialize the logger
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	// Initialize the config
	envConfigs := configs.InitConfigs()
	shutdownManager := shutdown.NewShutdownManager(logger)

	// Initialize the database connection pool
	dbConfig := db.DbConfig{Logger: logger, Ctx: ctx, DbUrl: envConfigs.Db.Url, ShutdownManager: shutdownManager}
	dbInstance := db.NewDb(dbConfig)
	pool,err := dbInstance.InitDB()

	if err != nil {
		logger.Error("failed to initialize the database connection pool",slog.Any("error", err))
		os.Exit(1)
		return nil,nil
	}
	
	movieRepository := db.NewMovieRepository(pool)
	createMovieUseCase := application.NewCreateMovieUsecase(movieRepository)
	getMovieUseCase := application.NewGetMovieUsecase(movieRepository)
	getMoviesUseCase := application.NewGetMoviesUsecase(movieRepository)
	updateMovieUseCase := application.NewUpdateMovieUsecase(movieRepository)
	deleteMovieUseCase := application.NewDeleteMovieUsecase(movieRepository)
	handlerConfig := api.HandlerConfig{
		Logger:              logger,
		CreateMovieUseCase:  createMovieUseCase,
		GetMovieUseCase:     getMovieUseCase,
		GetMoviesUseCase:    getMoviesUseCase,
		UpdateMovieUseCase:  updateMovieUseCase,
		DeleteMovieUseCase:  deleteMovieUseCase,	
	}
	handler := api.NewHandler(handlerConfig)
	
	return handler,pool
}

func createMovieRequest(t *testing.T, r *chi.Mux,payload *api.CreateMovieRequest) *http.Response {
	t.Helper()
	// Serialize payload
    body,err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	// Construct request
	req := httptest.NewRequest(http.MethodPost,"/api/v1/movies",bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")

	w := httptest.NewRecorder()
	// Make request
	r.ServeHTTP(w,req)
	return w.Result()
}

func createMovieAndGetID(t *testing.T, r *chi.Mux, payload *api.CreateMovieRequest) string {
	t.Helper()
	resp := createMovieRequest(t, r, payload)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	data, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
	}

	createdID, ok := data["id"].(string)
	if !ok {
		t.Fatalf("expected id to be a string, got %T", data["id"])
	}
	return createdID
}

func updateMovieRequest(t *testing.T, r *chi.Mux, id string, payload *api.UpdateMovieRequest) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/movies/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func getMovieRequest(t *testing.T, r *chi.Mux, id string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/"+id, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func getMoviesRequest(t *testing.T, r *chi.Mux) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func deleteMovieRequest(t *testing.T, r *chi.Mux, id string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/movies/"+id, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func resetMovieData(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if pool == nil {
		return
	}
	_, err := pool.Exec(context.Background(), `TRUNCATE TABLE movie_genres, movies RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("expected to reset movie data, got %v", err)
	}
	_, err = pool.Exec(context.Background(), `
		INSERT INTO genres (name) VALUES
			('ACTION'),
			('COMEDY'),
			('DRAMA'),
			('THRILLER'),
			('SCIENCE_FICTION'),
			('ADVENTURE'),
			('ROMANCE'),
			('MYSTERY'),
			('ANIMATION')
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("expected to seed genre data, got %v", err)
	}
}

func TestCreateMovieSuccess(t *testing.T) {
	handler,pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t,handler)
	payload := getCreateMovieSchema()

	resp := createMovieRequest(t,r,payload)
	
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code %d, got %d",http.StatusCreated,resp.StatusCode)
	}
	successResponse := decodeResponse[response.SuccessResponse](t,resp)
	if successResponse.State != "success" {
		t.Fatalf("expected state %s, got %s", "success",successResponse.State)
	}
	data, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
	}
	if data["title"] != payload.Title {
		t.Errorf("expected title %s, got %s", payload.Title,data["title"])
	}
}

func TestGetMovieSuccess(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	payload := getCreateMovieSchema(func(req *api.CreateMovieRequest) {
		req.Title = "get-movie-title"
	})
	createdID := createMovieAndGetID(t, r, payload)

	resp := getMovieRequest(t, r, createdID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	data, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
	}
	if data["id"] != createdID {
		t.Errorf("expected id %s, got %v", createdID, data["id"])
	}
	if data["title"] != payload.Title {
		t.Errorf("expected title %s, got %v", payload.Title, data["title"])
	}
}

func TestGetMovieNotFound(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	missingID := "11111111-1111-4111-8111-111111111111"

	resp := getMovieRequest(t, r, missingID)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != "movie_not_found" {
		t.Errorf("expected error code %s, got %s", "movie_not_found", errorResponse.Error.Code)
	}
}

func TestGetMoviesSuccess(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	payload := getCreateMovieSchema(func(req *api.CreateMovieRequest) {
		req.Title = "list-movie-title"
	})
	createdID := createMovieAndGetID(t, r, payload)

	resp := getMoviesRequest(t, r)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	data, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
	}
	movies, ok := data["movies"].([]any)
	if !ok {
		t.Fatalf("expected movies to be a list, got %T", data["movies"])
	}
	found := false
	for _, movie := range movies {
		movieMap, ok := movie.(map[string]any)
		if !ok {
			continue
		}
		if movieMap["id"] == createdID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected created movie %s to be present in the list", createdID)
	}
}

func TestUpdateMovieSuccess(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	payload := getCreateMovieSchema(func(req *api.CreateMovieRequest) {
		req.Title = "update-movie-title"
	})
	createdID := createMovieAndGetID(t, r, payload)

	updatedTitle := "updated-title"
	updatedDescription := "updated-description"
	updatedGenres := []string{"DRAMA"}
	updatedDurationMinutes := 90
	updatedRating := "R"
	updatedReleaseDate := time.Now().Add(48 * time.Hour)
	updatedDueDate := updatedReleaseDate.Add(24 * time.Hour)
	updatePayload := getUpdateMovieSchema(func(req *api.UpdateMovieRequest) {
		req.Title = &updatedTitle
		req.Description = &updatedDescription
		req.Genres = updatedGenres
		req.DurationMinutes = &updatedDurationMinutes
		req.Rating = &updatedRating
		req.ReleaseDate = &updatedReleaseDate
		req.DueDate = &updatedDueDate
	})

	resp := updateMovieRequest(t, r, createdID, updatePayload)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	data, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
	}
	if data["id"] != createdID {
		t.Errorf("expected id %s, got %v", createdID, data["id"])
	}
	if data["title"] != updatedTitle {
		t.Errorf("expected title %s, got %v", updatedTitle, data["title"])
	}
	if data["description"] != updatedDescription {
		t.Errorf("expected description %s, got %v", updatedDescription, data["description"])
	}
}

func TestUpdateMovieNotFound(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	missingID := "22222222-2222-4222-8222-222222222222"
	updatePayload := getUpdateMovieSchema(func(req *api.UpdateMovieRequest) {
		updatedTitle := "missing-title"
		req.Title = &updatedTitle
	})

	resp := updateMovieRequest(t, r, missingID, updatePayload)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != "movie_not_found" {
		t.Errorf("expected error code %s, got %s", "movie_not_found", errorResponse.Error.Code)
	}
}

func TestDeleteMovieSuccess(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	payload := getCreateMovieSchema(func(req *api.CreateMovieRequest) {
		req.Title = "delete-movie-title"
	})
	createdID := createMovieAndGetID(t, r, payload)

	resp := deleteMovieRequest(t, r, createdID)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	var deletedAt *time.Time
	err := pool.QueryRow(context.Background(), `SELECT deleted_at FROM movies WHERE id = $1`, createdID).Scan(&deletedAt)
	if err != nil {
		t.Fatalf("expected to query movie deletion state, got %v", err)
	}
	if deletedAt == nil {
		t.Errorf("expected movie %s to be soft deleted, but deleted_at was nil", createdID)
	}
}

func TestDeleteMovieNotFound(t *testing.T) {
	handler, pool := wireupHandler(t)
	resetMovieData(t, pool)
	t.Cleanup(func() {
		resetMovieData(t, pool)
	})
	r := setupRouter(t, handler)
	missingID := "33333333-3333-4333-8333-333333333333"

	resp := deleteMovieRequest(t, r, missingID)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != "movie_not_found" {
		t.Errorf("expected error code %s, got %s", "movie_not_found", errorResponse.Error.Code)
	}
}

