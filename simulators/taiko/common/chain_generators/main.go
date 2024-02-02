package chaingenerators

import (
	"github.com/ethereum/go-ethereum/core/types"
	el "github.com/taikoxyz/hive/simulators/taiko/common/config/execution"
)

type ChainGenerator interface {
	Generate(*el.ExecutionGenesis) ([]*types.Block, error)
}
