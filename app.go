package main

import (
	"context"

	"bifrost-client/internal/auth"
)

// App struct
type App struct {
	ctx  context.Context
	auth *auth.Service
}

// NewApp creates a new App application struct
func NewApp(authService *auth.Service) *App {
	return &App{
		auth: authService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


// CheckAndStartLogin checks authentication status and starts login if needed
func (a *App) CheckAndStartLogin() error {
	return a.auth.CheckAndStartLogin(a.ctx)
}
