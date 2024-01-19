package builder

import (
	"fmt"
	"lpm_builder/pkg/common"
	"lpm_builder/pkg/template"
	"os"
	"strings"
	"time"
)

func GenerateIndexPatch(ctx *BuilderCtx, build *template.Build) {
	tag := ""

	if build.Version.Tag != nil {
		tag = *build.Version.Tag
	}

	timestamp := time.Now().UTC().Unix()

	insertPart := "INSERT INTO repository (name, description, maintainer, source_repository, kind, tags, installed_size, license, v_major, v_minor, v_patch, v_tag, v_readable, mandatory_dependencies, suggested_dependencies, index_timestamp)"

	var mandatoryDependencies []string
	for _, val := range build.MandatoryDependencies.Runtime {
		pkg_with_version := fmt.Sprintf("%s@%s%s", val.Name, val.Version.Condition, val.Version.ReadableFormat)
		mandatoryDependencies = append(mandatoryDependencies, pkg_with_version)
	}

	var suggestedDependencies []string
	for _, val := range build.SuggestedDependencies.Runtime {
		pkg_with_version := fmt.Sprintf("%s@%s%s", val.Name, val.Version.Condition, val.Version.ReadableFormat)
		suggestedDependencies = append(suggestedDependencies, pkg_with_version)
	}

	valuesPart := fmt.Sprintf(`VALUES ("%s", "%s", "%s", "%s", "%s", "%s", %d, "%s", %d, %d, %d, "%s", "%s", "%s", "%s", %d);`,
		ctx.TemplateFields.Name,
		ctx.TemplateFields.Description,
		ctx.TemplateFields.Maintainer,
		ctx.TemplateFields.SourceRepository,
		ctx.TemplateFields.Kind,
		strings.Join(ctx.TemplateFields.Tags, ","),
		ctx.InstallSize,
		ctx.TemplateFields.License,
		build.Version.Major,
		build.Version.Minor,
		build.Version.Patch,
		tag,
		build.Version.ReadableFormat,
		strings.Join(mandatoryDependencies, ","),
		strings.Join(suggestedDependencies, ","),
		timestamp,
	)

	finalVersion := fmt.Sprintf("%s %s", insertPart, valuesPart)

	outputName := fmt.Sprintf("%d-%s.sql", timestamp, ctx.TemplateFields.Name)
	err := os.WriteFile(outputName, []byte(finalVersion), 0644)
	common.FailOnError(err)
}
