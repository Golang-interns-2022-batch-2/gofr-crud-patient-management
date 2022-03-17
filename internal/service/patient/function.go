package patient

func validatename(name string) bool {
	if name == "" {
		return name != ""
	}

	return true
}
