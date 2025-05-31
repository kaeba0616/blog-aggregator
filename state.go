package main

import (
	"github.com/kaeba0616/blog-aggregator/internal/config"
	"github.com/kaeba0616/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
