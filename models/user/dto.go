package user

type Register struct {
	UserType      string `json:"user_type" binding:"required"`
	Password      string `json:"password" binding:"required" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
	Name          string `json:"name" binding:"required" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
	Surname       string `json:"surname" binding:"required" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов."`
	Email         string `json:"email" binding:"required" valid:"email"`
	Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(7|18)~номер телефона от 7 до 18 цифр"`
	SocialNetwork string `json:"social_network"`
}

type Update struct {
	Name          string `json:"name" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
	Surname       string `json:"surname" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов."`
	Email         string `json:"email" valid:"email"`
	NewPassword   string `json:"new_password" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
	OldPassword   string `json:"old_password" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
	Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(4|18)~номер телефона от 4 до 18 цифр"`
	SocialNetwork string `json:"social_network"`
	Avatar        string `json:"avatar"`
}
