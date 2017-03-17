package kuddle

import (
	"path/filepath"

	rdef "go.polydawn.net/repeatr/api/def"
)

type FormulaLoader func(key string) (*rdef.Formula, error)

func FormulaLoaderForPath(basedir string) FormulaLoader {
	return func(key string) (*rdef.Formula, error) {
		_ = filepath.Clean(key)
		// ... TODO
		return nil, nil
	}
}

func FormulaLoaderConstant(frm *rdef.Formula) FormulaLoader {
	return func(_ string) (*rdef.Formula, error) {
		return frm.Clone(), nil
	}
}
