package main

import (
	"os"
	"pempal/resources"
)

// Path is the path of directories to search for resources.
var ResourcePath = resources.NewResourceFilePath(os.Getenv(ENV_PP_PATH))
