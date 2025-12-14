package models

// Group структура для группы пользователя
type Group struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// AdditionalInfo структура для дополнительной информации
type AdditionalInfo struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	DataType string `json:"dataType"`
}

// AccountInfoResponse структура ответа с полной информацией об аккаунте
type AccountInfoResponse struct {
	ID                         int              `json:"id"`
	Role                       string           `json:"role"`
	PositionID                 int              `json:"positionId"`
	PositionName               string           `json:"positionName"`
	DepartmentID               int              `json:"departmentId"`
	DepartmentName             string           `json:"departmentName"`
	FirstName                  string           `json:"firstName"`
	MiddleName                 string           `json:"middleName"`
	LastName                   string           `json:"lastName"`
	Email                      string           `json:"email"`
	Phone                      string           `json:"phone"`
	State                      string           `json:"state"`
	AvatarMinID                int              `json:"avatarMinId"`
	AvatarMidID                int              `json:"avatarMidId"`
	AvatarMaxID                int              `json:"avatarMaxId"`
	AvatarMinLink              string           `json:"avatarMinLink"`
	AvatarMidLink              string           `json:"avatarMidLink"`
	AvatarMaxLink              string           `json:"avatarMaxLink"`
	CommunicationMessangerLink string           `json:"communicationMessangerLink"`
	CommunicationMessangerIcon string           `json:"communicationMessangerIcon"`
	Groups                     []Group          `json:"groups"`
	AdditionalInfo             []AdditionalInfo `json:"additionalInfo"`
	StatusCode                 int              `json:"-"`
	Message                    string           `json:"-"`
}
