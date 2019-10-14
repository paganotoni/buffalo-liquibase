package migrator

type Migrator interface {
	Prepare() error
	Up() error
	Down() error
}
