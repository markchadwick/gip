package main

func get(imprt string) error {
	path, err := OpenPath()
	if err != nil {
		return err
	}

	pkg, err := PackageFromImport(imprt)
	if err != nil {
		return err
	}

	return path.Install(pkg)
}
