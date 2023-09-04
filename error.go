package main

func (r Report) validateReport() []FormError {
	var formErrors []FormError

	if r.Title == "" {
		formErrors = append(formErrors, FormError{
			Field:   "title",
			Message: "Title cannot be empty",
		})
	}

	if r.Type == "" {
		formErrors = append(formErrors, FormError{
			Field:   "type",
			Message: "Type cannot be empty",
		})
	}

	if r.Description == "" {
		formErrors = append(formErrors, FormError{
			Field:   "description",
			Message: "Description cannot be empty",
		})
	}

	if r.Image == "" {
		formErrors = append(formErrors, FormError{
			Field:   "image",
			Message: "Image cannot be empty",
		})
	}

	return formErrors
}

func (u User) validateUser() []FormError {
	unameAji := envString("UNAME_AJI", aji)
	pwAji := envString("PW_AJI", migrationOnPw)

	var formErrors []FormError

	if u.Name == "" {
		formErrors = append(formErrors, FormError{
			Field:   "name",
			Message: "Name cannot be empty",
		})
	}

	if u.Password == "" {
		formErrors = append(formErrors, FormError{
			Field:   "password",
			Message: "Password cannot be empty",
		})
	}

	if u.Name != unameAji {
		formErrors = append(formErrors, FormError{
			Field:   "name",
			Message: "Name is not valid",
		})
	}
	if u.Password != pwAji {
		formErrors = append(formErrors, FormError{
			Field:   "password",
			Message: "Password is not valid",
		})
	}

	return formErrors
}
