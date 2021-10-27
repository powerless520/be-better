package model

import (
	"be-better/utils/idGenerator"
)

type IdGenerators struct {
	UserIdGenerator    *idGenerator.IdGenerator
	SessionIdGenerator *idGenerator.IdGenerator
}
