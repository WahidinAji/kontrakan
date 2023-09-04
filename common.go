package main

import "gopkg.in/guregu/null.v4"

type User struct {
	Name      string      `json:"name"`
	Password  string      `json:"password"`
	CreatedAt null.String `json:"created_at"`
}

type Report struct {
	Id          uint64      `json:"id"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	CreatedAt   null.String `json:"created_at"`
	UpdatedAt   null.String `json:"updated_at"`
}

type Response struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	FormErrors []FormError `json:"form_errors"`
}

type FormError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
