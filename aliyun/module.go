package aliyun

import (
	"embed"

	"github.com/starter-go/application"
	"github.com/starter-go/bucket-drivers/aliyun/gen/main4aliyun"
	"github.com/starter-go/bucket-drivers/aliyun/gen/test4aliyun"
	"github.com/starter-go/buckets/modules/buckets"
	"github.com/starter-go/starter"
)

////////////////////////////////////////////////////////////////////////////////

const (
	theModuleName     = "github.com/starter-go/bucket-drivers/aliyun"
	theModuleVersion  = "v0.0.1"
	theModuleRevision = 3
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

	mb.Components(main4aliyun.ExportComponents)

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

	mb.Components(test4aliyun.ExportComponents)

	mb.Depend(Module())

	return mb.Create()
}

////////////////////////////////////////////////////////////////////////////////
// EOF
