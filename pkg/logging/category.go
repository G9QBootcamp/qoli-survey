package logging

type Category string
type SubCategory string
type ExtraKey string

const (
	General         Category = "General"
	Internal        Category = "Internal"
	Database        Category = "Database"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Io              Category = "Io"
)

const (
	// General
	Startup         SubCategory = "Startup"
	ExternalService SubCategory = "ExternalService"

	// Database
	Migration       SubCategory = "Migration"
	Select          SubCategory = "Select"
	DatabaseTimeout SubCategory = "DatabaseTimeout"
	Rollback        SubCategory = "Rollback"
	Update          SubCategory = "Update"
	Delete          SubCategory = "Delete"
	Insert          SubCategory = "Insert"

	//validation
	QuestionnaireValidationFailed SubCategory = "QuestionnaireValidationFailed"

	// Internal
	Api                 SubCategory = "Api"
	HashPassword        SubCategory = "HashPassword"
	UserNotAuthorized   SubCategory = "UserNotAuthorized"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"
	FailedToCreateUser  SubCategory = "FailedToCreateUser"
	FailedToGetUsers    SubCategory = "FailedToGetUsers"

	RecoverError SubCategory = "Recover Error"
)

const (
	AppName       ExtraKey = "AppName"
	LoggerName    ExtraKey = "Logger"
	Service       ExtraKey = "Service"
	UserId        ExtraKey = "UserId"
	TraceId       ExtraKey = "TraceId"
	ClientIp      ExtraKey = "ClientIp"
	HostIp        ExtraKey = "HostIp"
	Method        ExtraKey = "Method"
	StatusCode    ExtraKey = "StatusCode"
	BodySize      ExtraKey = "BodySize"
	Path          ExtraKey = "Path"
	Latency       ExtraKey = "Latency"
	RequestBody   ExtraKey = "RequestBody"
	RequestHeader ExtraKey = "RequestHeader"
	ResponseBody  ExtraKey = "ResponseBody"
	ErrorMessage  ExtraKey = "ErrorMessage"
)
