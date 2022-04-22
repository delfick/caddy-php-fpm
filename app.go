package caddyphpfpm

import (
	"github.com/caddyserver/caddy/v2"
)

// Interface guards
var (
	_ caddy.App          = (*App)(nil)
	_ caddy.Module       = (*App)(nil)
	_ caddy.Provisioner  = (*App)(nil)
	_ caddy.CleanerUpper = (*App)(nil)
)

type App struct {
	php *PHP
}

// NewApp returns a new App object
func NewApp() *App {
	return &App{php: NewPHP()}
}

// CaddyModule implements caddy.Module
func (a *App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "php_fpm",
		New: func() caddy.Module { return a },
	}
}

func (a *App) Provision(ctx caddy.Context) error {
	go a.php.Run()
	return a.php.wait()
}

// Start implements caddy.App
func (a *App) Start() error {
	return nil
}

func (a *App) Cleanup() error {
	a.php.Stop()
	return nil
}

// Stop implements caddy.App
func (a *App) Stop() error {
	return nil
}
