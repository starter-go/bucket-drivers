package tencent

import (
	"embed"

	"github.com/starter-go/application"
	"github.com/starter-go/bucket-drivers/tencent/gen/main4tencent"
	"github.com/starter-go/bucket-drivers/tencent/gen/test4tencent"

	"github.com/starter-go/buckets/modules/buckets"
	"github.com/starter-go/starter"
)

////////////////////////////////////////////////////////////////////////////////

const (
	theModuleName     = "github.com/starter-go/bucket-drivers/tencent"
	theModuleVersion  = "v0.0.0"
	theModuleRevision = 0
)

////////////////////////////////////////////////////////////////////////////////

const (
	theMainModuleResPath = "src/main/resources"
	theTestModuleResPath = "src/test/resources"
)

//go:embed "src/main/resources"
var theMainModuleResFS embed.FS

//go:embed "src/test/resources"
var theTestModuleResFS embed.FS

////////////////////////////////////////////////////////////////////////////////

func Module() application.Module {
	mb := new(application.ModuleBuilder)

	mb.Name(theModuleName + "#main")
	mb.Version(theModuleVersion)
	mb.Revision(theModuleRevision)

	mb.EmbedResources(theMainModuleResFS, theMainModuleResPath)

	mb.Components(main4tencent.ExportComponents)

	mb.Depend(starter.Module())
	mb.Depend(buckets.ModuleLib())

	return mb.Create()
}

func ModuleForTest() application.Module {
	mb := new(application.ModuleBuilder)

	mb.Name(theModuleName + "#test")
	mb.Version(theModuleVersion)
	mb.Revision(theModuleRevision)

	mb.EmbedResources(theTestModuleResFS, theTestModuleResPath)

	mb.Components(test4tencent.ExportComponents)

	mb.Depend(Module())

	return mb.Create()
}

////////////////////////////////////////////////////////////////////////////////
// EOF
