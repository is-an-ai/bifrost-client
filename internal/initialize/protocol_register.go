package initialize

import (
	"context"
	"log"
	"strings"

	"bifrost-client/internal/auth"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// RegisterProtocolHandlers registers protocol handlers for the application
func RegisterProtocolHandlers(authService *auth.Service) (mac.Options, options.SingleInstanceLock) {
	return mac.Options{
			OnUrlOpen: func(url string) {
				handleCallbackURL(context.Background(), authService, url)
			},
		}, options.SingleInstanceLock{
			UniqueId: "bifrost-client-auth",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				if len(secondInstanceData.Args) > 0 {
					handleCallbackURL(context.Background(), authService, secondInstanceData.Args[0])
				}
			},
		}
}

// handleCallbackURL processes the callback URL from the protocol handler
func handleCallbackURL(ctx context.Context, authService *auth.Service, url string) {
	if strings.HasPrefix(url, "bifrost://auth/callback") {
		if err := authService.HandleCallback(ctx, url); err != nil {
			log.Printf("Failed to handle callback: %v", err)
		}
	}
}
