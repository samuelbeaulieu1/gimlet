package responses

import "github.com/samuelbeaulieu1/gimlet/lang"

type Response int
type ResponseHandler interface {
	String(Response) string
	Error(Response) Error
}

type ResponseMapper struct{}

const (
	ERR_RECORD_NOT_FOUND Response = iota
	ERR_GET_RECORDS
	ERR_GET_RECORD
	ERR_UPDATE_RECORD
	ERR_CREATE_RECORD
	ERR_DELETE_RECORD

	ERR_EMPTY_ID
	ERR_INVALID_ID

	SUCC_UPDATE_RECORD
	SUCC_DELETE_RECORD
)

var fr = []string{
	"Record inexistant",
	"Une erreur est survenue en récupérant la liste de records",
	"Une erreur est survenue en récupérant le record",
	"Une erreur est survenue en modifiant le record",
	"Une erreur est survenue en créant le record",
	"Une erreur est survenue en supprimant le record",

	"L'identifiant est obligatoire",
	"L'identifiant est invalide",

	"Le record a été modifié",
	"Le record record a été supprimé",
}

var en = []string{
	"Record not found",
	"An unexpected error has occured while retrieving the records",
	"An unexpected error has occured while retrieving the record",
	"An unexpected error has occured while trying to update the record",
	"An unexpected error has occured while trying to create the record",
	"An unexpected error has occured while trying to delete the record",

	"The identifier is missing",
	"The identifier is invalid",

	"The record was updated",
	"The record was deleted",
}

var messages = [][]string{
	fr,
	en,
}

func (response Response) String() string {
	return messages[lang.Get()][response]
}

func (response Response) Error() Error {
	return NewError(response.String())
}

func (mapper *ResponseMapper) String(response Response) string {
	return response.String()
}

func (mapper *ResponseMapper) Error(response Response) Error {
	return response.Error()
}
