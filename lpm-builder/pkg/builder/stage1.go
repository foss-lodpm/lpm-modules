package builder

import (
	"lpm_builder/pkg/common"
	"path/filepath"
)

const (
	PreInstall  = "pre_install"
	PostInstall = "post_install"

	PreDelete  = "pre_delete"
	PostDelete = "post_delete"

	PreDowngrade  = "pre_downgrade"
	PostDowngrade = "post_downgrade"

	PreUpgrade  = "pre_upgrade"
	PostUpgrade = "post_upgrade"
)

func CopyProvidedStage1Scripts(ctx *BuilderCtx) {
	scripts := []string{
		PreInstall,
		PostInstall,
		PreDelete,
		PostDelete,
		PreDowngrade,
		PostDowngrade,
		PreUpgrade,
		PostUpgrade,
	}

	for _, script := range scripts {
		srcPath := filepath.Join(ctx.Stage1ScriptsDir, script)
		destPath := filepath.Join(ctx.TmpScriptsDir, script)
		err := common.CopyIfExists(srcPath, destPath)
		common.FailOnError(err)
	}
}
