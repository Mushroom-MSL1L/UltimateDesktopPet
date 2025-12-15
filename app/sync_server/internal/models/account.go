package models

type Account struct {
	Username     string `json:"username" example:"admin"`
	Password string `json:"password" example:"1234"`
}

func GetTheAccount() Account {
	return Account{
		Username: "admin",
		Password: "1234",
	}
}

func ValidateUsername(username string) bool {
	theAccount := GetTheAccount()
	return username == theAccount.Username
}

func ValidateAccount(account Account) bool {
	theAccount := GetTheAccount()
	return (account.Username == theAccount.Username) && (account.Password == theAccount.Password)
}