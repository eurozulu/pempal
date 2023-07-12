package utils

import "bytes"

type CompoundErrors []error

func (errs CompoundErrors) Error() string {
	buf := bytes.NewBuffer(nil)
	for i, e := range errs {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(e.Error())
	}
	return buf.String()
}
