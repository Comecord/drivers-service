package service_errors

const (
	// Token
	UnExpectedError = "Expected error"
	ClaimsNotFound  = "Claims not found"
	TokenRequired   = "Требуеться токен для авторизации в системе"
	TokenBearer     = "Требуеться токен для авторизации в системе"
	TokenExpired    = "Токен устарел"
	TokenInvalid    = "Токен недействителен"

	// OTP
	OtpExists   = "Otp exists"
	OtpUsed     = "Otp used"
	OtpNotValid = "Otp invalid"

	// User
	EmailExists       = "Указанный email уже зарегистрирован в системе"
	EmailNotExists    = "Email не зарегистрирован в системе"
	UsernameExists    = "Пользователь с таким именем уже зарегистрирован в системе"
	UsernameNotExists = "Пользователь с таким именем не зарегистрирован в системе"
	PermissionDenied  = "У вас нет прав для этого действия"

	// DB
	RecordNotFound = "Запись не найдена"

	// Role
	RoleExists = "Роль с таким названием уже существует"
)
