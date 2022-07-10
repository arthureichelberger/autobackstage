package model

const RefDoesNotExistError string = "Reference does not exist"

type ResponseError struct {
	Message string `json:"message"`
}
