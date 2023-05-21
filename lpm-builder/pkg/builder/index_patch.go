package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
	"os"
	"time"
)

func GenerateIndexPatch(ctx *BuilderCtx) {
	tag := ""

	if ctx.TemplateFields.Version.Tag != nil {
		tag = *ctx.TemplateFields.Version.Tag
	}

	timestamp := time.Now().UTC().Unix()

	insertPart := "INSERT INTO repository (name, description, maintainer, kind, installed_size, license, v_major, v_minor, v_patch, v_tag, v_readable, index_timestamp)"

	valuesPart := fmt.Sprintf(`VALUES ("%s", "%s", "%s", "%s", %d, "%s", %d, %d, %d, "%s", "%s", %d);`,
		ctx.TemplateFields.Name,
		ctx.TemplateFields.Description,
		ctx.TemplateFields.Maintainer,
		ctx.TemplateFields.Kind,
		ctx.InstallSize,
		ctx.TemplateFields.License,
		ctx.TemplateFields.Version.Major,
		ctx.TemplateFields.Version.Minor,
		ctx.TemplateFields.Version.Patch,
		tag,
		ctx.TemplateFields.Version.ReadableFormat,
		timestamp,
	)

	finalVersion := insertPart + " " + valuesPart

	outputName := fmt.Sprintf("%d-%s.sql", timestamp, ctx.TemplateFields.Name)
	err := os.WriteFile(outputName, []byte(finalVersion), 0644)
	common.FailOnError(err)
}
