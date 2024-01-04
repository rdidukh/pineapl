package ast

import (
	"log"

	"github.com/llir/llvm/ir/types"
)

type Type string

func (typ Type) toIrType() types.Type {
	switch typ {
	case "Bool":
		return types.I1
	case "Int":
		return types.I32
	case "Float":
		return types.Float
	}

	log.Panicf("unsupported type: %s", typ)
	return types.Void
}
