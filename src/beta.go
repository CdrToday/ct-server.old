package main

func betaAccount(mail string) bool {
	accounts := []string{
		"test@cdr.today",
		"john@cdr.today",
		"mercury@cdr.today",
	}

	for _, i := range accounts {
		if mail == i {
			return true
		}
	}

	return false
}
