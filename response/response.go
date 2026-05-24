package response

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

func Success(data interface{}) Response {
    return Response{
        Success: true,
        Data:    data,
    }
}

func SuccessWithMessage(data interface{}, message string) Response {
    return Response{
        Success: true,
        Data:    data,
        Message: message,
    }
}

func Error(err string) Response {
    return Response{
        Success: false,
        Error:   err,
    }
}

func ErrorResponse(err string) Response {
    return Error(err)
}

func ErrorWithMessage(err string, message string) Response {
    return Response{
        Success: false,
        Error:   err,
        Message: message,
    }
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        // Log error but don't panic
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func WriteError(w http.ResponseWriter, statusCode int, errMsg string) {
    WriteJSON(w, statusCode, ErrorResponse(errMsg))
}

func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
    WriteJSON(w, statusCode, Success(data))
}