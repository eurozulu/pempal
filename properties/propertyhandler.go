package properties

import "github.com/eurozulu/pempal/templates"

type PropertyHandler interface {
	HandleError(err string, template templates.Template) error
}

func PropertyHandlerForName(name string) (PropertyHandler, error) {

}
