package middleware

import (
    "github.com/rs/cors"
)

func EnableCORS() *cors.Cors {
    return cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
        Debug:            false,
    })
}