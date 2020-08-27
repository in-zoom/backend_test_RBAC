package data

type Message struct {
	OkMessage string `json:"okMessage"`
	Status    int    `json:"status"`
}
type ErrMessage struct {
	Message string `json:"message"`
}

type RegisterUser struct {
	Name     string `json:"name"`
	Role     string `json:"email"`
	Password string `json:"password"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	AdminName string = "Admin"
	Role      string = "admin"
	AdminPass string = "123456789"
	UserRole  string = "user"
)
