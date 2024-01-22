package http_response

type App string
type AppError string

const (
	MODE_DEV                  App      = "dev"
	MODE_PROD                 App      = "prod"
	SUCCESS                   string   = "SUCCESS"
	NOT_FOUND                 AppError = "NOT_FOUND"
	USER_ID_NOT_FOUND         AppError = "USER_ID_NOT_FOUND"
	USER_DATA_NOT_FOUND       AppError = "USER_DATA_NOT_FOUND"
	PASSWORD_NOT_MATCH        AppError = "PASSWORD_NOT_MATCH"
	ACCESS_DENIED             AppError = "ACCESS_DENIED"
	API_ACCESS_DENIED         AppError = "API_ACCESS_DENIED"
	ID_NOT_EXITST             AppError = "ID_NOT_EXITS"
	UNABLE_TO_PROCESS_REQUEST AppError = "UNABLE_TO_PROCESS_REQUEST"
	BAD_REQUEST               AppError = "BAD_REQUEST"
)

func (e AppError) Error() string {
	return string(e)
}
