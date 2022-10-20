package file

import "github.com/elvuel/temprender/transport"

const (
	FileCreatorTransporterKind = "file::creator"
	FileCreatorAliasExt        = ".alias_newer"

	FileInjectorTransporterKind = "file::injector"
	FileInjectorAliasExt        = ".alias_changed"

	FileDestroyerTransporterKind = "file::destroyer"
)

func init() {
	transport.RegisterTransporter(&transport.TransporterManifest{
		Kind:    FileCreatorTransporterKind,
		NewFunc: NewCreatorRegister,
	})

	transport.RegisterTransporter(&transport.TransporterManifest{
		Kind:    FileInjectorTransporterKind,
		NewFunc: NewInjectorRegister,
	})

	transport.RegisterTransporter(&transport.TransporterManifest{
		Kind:    FileDestroyerTransporterKind,
		NewFunc: NewDestroyerRegister,
	})
}
