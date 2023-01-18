package shared

import "os"

// リダイレクトURLです
var RedirectURL = os.Getenv("FE_ROOT_URL")
