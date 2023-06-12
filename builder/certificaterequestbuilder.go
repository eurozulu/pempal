package builder

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
)

type certificateRequestBuilder struct {
	dto  model.CertificateRequestDTO
	keys keys.Keys
}

func (crb certificateRequestBuilder) ApplyTemplate(tp ...templates.Template) error {
	//TODO implement me
	panic("implement me")
}

func (crb certificateRequestBuilder) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (crb certificateRequestBuilder) Build() (model.Resource, error) {
	//TODO implement me
	panic("implement me")
}

func (crb *certificateRequestBuilder) UnmarshalArguments(args []string) ([]string, error) {
	remain, err := argdecoder.ApplyArguments(args, &crb.dto)
	if err != nil {
		return nil, err
	}
	if len(remain) > 0 {
		return nil, fmt.Errorf("unknown flags: %v", remain)
	}
	return nil, nil
}
