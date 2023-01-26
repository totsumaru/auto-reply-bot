package redirect

import "os"

// リダイレクトURLです
func CreateRedirectURL() string {
	return os.Getenv("FE_ROOT_URL")
}
