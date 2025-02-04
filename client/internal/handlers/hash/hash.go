package hash

import (
	"github.com/mrkucher83/hash-service/client/internal/godb"
)

type ReqBody struct {
	Params []string
}

type Repo struct {
	storage *godb.Instance
}

func NewRepo(r *godb.Instance) *Repo {
	return &Repo{storage: r}
}
