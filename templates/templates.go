package templates

import (
	"os"
)

const ENV_TEMPLATE_PATH = "TEMPLATE_PATH"

var TemplatePath = NewTemplatePath(os.ExpandEnv(os.Getenv(ENV_TEMPLATE_PATH)))
