package cmd_logout

import (
	"github.com/reduan2660/swapenv/internal/api"
)

func Logout() error {
	return api.DeleteCredentials()
}
