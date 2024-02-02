package xsql

type MigrationData interface {
	Exist() (bool, error)
	Create() error
}

func MigrateData(dos []MigrationData) error {
	for _, do := range dos {
		ext, err := do.Exist()
		if err != nil {
			return err
		}
		if ext {
			continue
		}
		err = do.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
