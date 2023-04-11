package schemas

type UserDetailResponse struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	IsSuperuser bool   `json:"is_superuser"`
}

type UserPaginateResponse struct {
	Counts    int                  `json:"counts"`
	PageCount int                  `json:"page_count"`
	PageSize  int                  `json:"page_size"`
	Page      int                  `json:"page"`
	Results   []UserDetailResponse `json:"results"`
}

type UserCreateRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsActive    bool   `json:"is_active"`
	IsSuperuser bool   `json:"is_superuser"`
}

type UserCreateResponse struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	IsSuperuser bool   `json:"is_superuser"`
}

type UserUpdateRequest struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Password    *string `json:"password"`
	IsActive    bool    `json:"is_active"`
	IsSuperuser bool    `json:"is_superuser"`
}

type UserUpdateResponse struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	IsSuperuser bool   `json:"is_superuser"`
}
