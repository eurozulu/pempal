package resourceio

import "pempal/model"

type ResourceLocation struct {
	Path      string
	Resources []model.PEMResource
}
