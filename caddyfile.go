package caddyphpfpm

import (
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("php_fpm", parseOptions)
}

// parseOptions configures our options
// Syntax:
//  php_fpm {
//		sock_location fpm.sock
//      start_timeout 10s
//  }
func parseOptions(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := NewApp()

	// consume the option name
	if !d.Next() {
		return nil, d.ArgErr()
	}

	// handle any options
	for d.NextBlock(0) {
		switch d.Val() {
		case "cmd":
			app.php.command = d.RemainingArgs()

			if len(app.php.command) == 0 {
				return nil, d.ArgErr()
			}
		case "sock_location":
			if !d.Args(&app.php.sockLocation) {
				return nil, d.ArgErr()
			}
		case "start_timeout":
			var dur string
			if !d.Args(&dur) {
				return nil, d.ArgErr()
			}

			timeout, err := time.ParseDuration(dur)
			if err != nil {
				return nil, err
			}
			app.php.startTimeout = timeout
		}
	}

	caddy.RegisterModule(app)

	// tell Caddyfile adapter that this is the JSON for an app
	return httpcaddyfile.App{
		Name:  "php_fpm",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}
