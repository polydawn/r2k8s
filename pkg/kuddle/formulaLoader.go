package kuddle

import (
	"path/filepath"

	"go.polydawn.net/meep"
	rdef "go.polydawn.net/repeatr/api/def"
	rhitch "go.polydawn.net/repeatr/api/hitch"
)

type FormulaLoader func(key string) (*rdef.Formula, error)

func FormulaLoaderForPath(basedir string) FormulaLoader {
	return func(key string) (frm *rdef.Formula, err error) {
		pth := filepath.Clean(key)
		err = meep.RecoverPanics(func() {
			frm = rhitch.LoadFormulaFromFile(pth)
		})
		return
	}
}

func FormulaLoaderConstant(frm *rdef.Formula) FormulaLoader {
	return func(_ string) (*rdef.Formula, error) {
		return frm.Clone(), nil
	}
}
