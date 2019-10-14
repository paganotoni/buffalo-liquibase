package migrator

type Migrator interface {
	Prepare() error
	Up() error
	Down(int) error
}
