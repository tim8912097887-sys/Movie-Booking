package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/json"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/middleware"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/validation"
)

type Handler struct {
	logger              *slog.Logger
	createMovieUseCase  CreateMovieUsecase
	getMovieUseCase     GetMovieUsecase
	getMoviesUseCase    GetMoviesUsecase
	updateMovieUseCase  UpdateMovieUsecase
	deleteMovieUseCase  DeleteMovieUsecase
}

type HandlerConfig struct {
	Logger              *slog.Logger
	CreateMovieUseCase  CreateMovieUsecase
	GetMovieUseCase     GetMovieUsecase
	GetMoviesUseCase    GetMoviesUsecase
	UpdateMovieUseCase  UpdateMovieUsecase
	DeleteMovieUseCase  DeleteMovieUsecase
}

func NewHandler(handlerConfig HandlerConfig) *Handler {
	return &Handler{
		logger:             handlerConfig.Logger,
		createMovieUseCase: handlerConfig.CreateMovieUseCase,
		getMovieUseCase:    handlerConfig.GetMovieUseCase,
		getMoviesUseCase:   handlerConfig.GetMoviesUseCase,
		updateMovieUseCase: handlerConfig.UpdateMovieUseCase,
		deleteMovieUseCase: handlerConfig.DeleteMovieUseCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
   r.Post("/", h.CreateMovie)
   r.Get("/", h.GetMovies)
   r.Get("/{id}", h.GetMovie)
   r.Put("/{id}", h.UpdateMovie)
   r.Delete("/{id}", h.DeleteMovie)
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
    result, err := validation.ValidateRequestBody[CreateMovieRequest](r)
	if err != nil {
        h.logger.Error("failed to validate request body",slog.Any("error", err))
	    err = json.ErrorJson(w, http.StatusBadRequest, "validation_error", err.Error())
	    if err != nil {
			h.logger.Error("failed to write error response",slog.Any("error", err))
		}
		return
	}
    
    createMovieCommand := application.CreateMovieCommand{
		Title:           result.Title,
		Description:     result.Description,
		Genres:          result.Genres,
		DurationMinutes: result.DurationMinutes,
		Rating:          result.Rating,
		ReleaseDate:     result.ReleaseDate,
		DueDate:         result.DueDate,
	}
    
	createdMovie, err := h.createMovieUseCase.Execute(r.Context(), createMovieCommand)
	var responseError error
	if err != nil {
		h.logger.Error("failed to create movie",slog.Any("error", err))
		responseError = middleware.ErrorHandler(w,err)
		if responseError != nil {
			h.logger.Error("failed to write error response",slog.Any("error", responseError))
		}
		return
	}

	movieResponse := CreateMovieResponse{
		ID:              createdMovie.ID,
		Title:           createdMovie.Title,
		Description:     createdMovie.Description,
		Genres:          createdMovie.Genres,
		DurationMinutes: createdMovie.DurationMinutes,
		Rating:          createdMovie.Rating,
		ReleaseDate:     createdMovie.ReleaseDate,
		DueDate:         createdMovie.DueDate,
	}
	responseError = json.SuccessJson(w, http.StatusCreated, movieResponse)
	if responseError != nil {
		h.logger.Error("failed to write success response",slog.Any("error", responseError))
	}
}

func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {

	uuidRequest := UUIDRequest{
         ID: r.PathValue("id"),
	}
	result,err := validation.Validate(uuidRequest)
	if err != nil {
		h.logger.Error("failed to validate request body",slog.Any("error", err))
		err = json.ErrorJson(w, http.StatusBadRequest, "validation_error", err.Error())
		if err != nil {
			h.logger.Error("failed to write error response",slog.Any("error", err))
		}
		return
	}

	movie, err := h.getMovieUseCase.Execute(r.Context(), result.ID)
    var responseError error
	if err != nil {
		h.logger.Error("failed to get movie",slog.Any("error", err))
		responseError = middleware.ErrorHandler(w,err)
		if responseError != nil {
			h.logger.Error("failed to write error response",slog.Any("error", responseError))
		}
		return
	}

	movieResponse := GetMovieResponse{
		ID:              movie.ID,
		Title:           movie.Title,
		Description:     movie.Description,
		Genres:          movie.Genres,
		DurationMinutes: movie.DurationMinutes,
		Rating:          movie.Rating,
		ReleaseDate:     movie.ReleaseDate,
		DueDate:         movie.DueDate,
	}
	responseError = json.SuccessJson(w, http.StatusOK, movieResponse)
	if responseError != nil {
		h.logger.Error("failed to write success response",slog.Any("error", responseError))
	}
	
}


func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.getMoviesUseCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("failed to get movies", slog.Any("error", err))
		err = json.ErrorJson(w, http.StatusInternalServerError, "get_movies_error", err.Error())
		if err != nil {
			h.logger.Error("failed to write error response", slog.Any("error", err))
		}
		return
	}

	var moviesResponse []GetMovieResponse
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, GetMovieResponse{
			ID:              movie.ID,
			Title:           movie.Title,
			Description:     movie.Description,
			Genres:          movie.Genres,
			DurationMinutes: movie.DurationMinutes,
			Rating:          movie.Rating,
			ReleaseDate:     movie.ReleaseDate,
			DueDate:         movie.DueDate,
		})
	}

	response := GetMoviesResponse{
		Movies: moviesResponse,
	}
	var responseError error
	responseError = json.SuccessJson(w, http.StatusOK, response)
	if responseError != nil {
		h.logger.Error("failed to write success response", slog.Any("error", responseError))
	}
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	result, err := validation.ValidateRequestBody[UpdateMovieRequest](r)
	if err != nil {
		h.logger.Error("failed to validate request body", slog.Any("error", err))
		err = json.ErrorJson(w, http.StatusBadRequest, "validation_error", err.Error())
		if err != nil {
			h.logger.Error("failed to write error response", slog.Any("error", err))
		}
		return
	}

	uuidRequest := UUIDRequest{
		ID: r.PathValue("id"),
	}
	validationResult, err := validation.Validate(uuidRequest)
	if err != nil {
		h.logger.Error("failed to validate request body", slog.Any("error", err))
		err = json.ErrorJson(w, http.StatusBadRequest, "validation_error", err.Error())
		if err != nil {
			h.logger.Error("failed to write error response", slog.Any("error", err))
		}
		return
	}

	updateMovieCommand := application.UpdateMovieCommand{
		ID:              validationResult.ID,
		Title:           result.Title,
		Description:     result.Description,
		Genres:          result.Genres,
		DurationMinutes: result.DurationMinutes,
		Rating:          result.Rating,
		ReleaseDate:     result.ReleaseDate,
		DueDate:         result.DueDate,
	}

	updatedMovie, err := h.updateMovieUseCase.Execute(r.Context(), updateMovieCommand)
	var responseError error
	if err != nil {
		h.logger.Error("failed to update movie", slog.Any("error", err))
		responseError = middleware.ErrorHandler(w, err)
		if responseError != nil {
			h.logger.Error("failed to write error response", slog.Any("error", responseError))
		}
		return
	}

	movieResponse := UpdateMovieResponse{
		ID:              updatedMovie.ID,
		Title:           updatedMovie.Title,
		Description:     updatedMovie.Description,
		Genres:          updatedMovie.Genres,
		DurationMinutes: updatedMovie.DurationMinutes,
		Rating:          updatedMovie.Rating,
		ReleaseDate:     updatedMovie.ReleaseDate,
		DueDate:         updatedMovie.DueDate,
	}
	responseError = json.SuccessJson(w, http.StatusOK, movieResponse)
	if responseError != nil {
		h.logger.Error("failed to write success response", slog.Any("error", responseError))
	}
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	uuidRequest := UUIDRequest{
		ID: r.PathValue("id"),
	}
	result, err := validation.Validate(uuidRequest)
	if err != nil {
		h.logger.Error("failed to validate request body", slog.Any("error", err))
		err = json.ErrorJson(w, http.StatusBadRequest, "validation_error", err.Error())
		if err != nil {
			h.logger.Error("failed to write error response", slog.Any("error", err))
		}
		return
	}

	deleteMovieCommand := application.DeleteMovieCommand{
		ID: result.ID,
	}

	_, err = h.deleteMovieUseCase.Execute(r.Context(), deleteMovieCommand)
	var responseError error
	if err != nil {
		h.logger.Error("failed to delete movie", slog.Any("error", err))
		responseError = middleware.ErrorHandler(w, err)
		if responseError != nil {
			h.logger.Error("failed to write error response", slog.Any("error", responseError))
		}
		return
	}

	responseError = json.SuccessJson(w, http.StatusNoContent, nil)
	if responseError != nil {
		h.logger.Error("failed to write success response", slog.Any("error", responseError))
	}
}