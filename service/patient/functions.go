package patient

func IsIDValid(id int) bool {
	return id > 0
}
func IsNameValid(name string) bool {
	return len(name) > 0
}
func IsPhoneValid(phone string) bool {
	return len(phone) == 13 && phone[:3] == "+91"
}
