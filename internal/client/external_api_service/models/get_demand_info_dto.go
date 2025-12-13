package models

// StatusMeta структура для метаданных статуса
type StatusMeta struct {
	IconPath  string `json:"iconPath"`
	IconColor string `json:"iconColor"`
}

// PreviewData структура для данных предпросмотра
type PreviewData struct {
	ThemeName               string            `json:"themeName"`
	ShowStatusIcons         bool              `json:"showStatusIcons"`
	ShowRequestCreationDate bool              `json:"showRequestCreationDate"`
	ShowRequestID           bool              `json:"showRequestId"`
	AccountData             map[string]string `json:"accountData"`
	FormData                map[string]string `json:"formData"`
}

// FormCellAnswer структура для ответов в ячейках формы
type FormCellAnswer struct {
	Group       int      `json:"group"`
	OrderNumber int      `json:"orderNumber"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Answer      []string `json:"answer"`
}

// DemandInfoResponse структура ответа с полной информацией о заявке
type DemandInfoResponse struct {
	Category                        string             `json:"category"`
	Subject                         string             `json:"subject"`
	SectionName                     string             `json:"sectionName"`
	ID                              int                `json:"id"`
	CreateDate                      string             `json:"createDate"`
	Status                          string             `json:"status"`
	IsChecked                       bool               `json:"isChecked"`
	AccountFirstName                string             `json:"accountFirstName"`
	AccountLastName                 string             `json:"accountLastName"`
	AccountMiddleName               string             `json:"accountMiddleName"`
	PriorityID                      int                `json:"priorityId"`
	PriorityName                    string             `json:"priorityName"`
	PriorityColor                   string             `json:"priorityColor"`
	StatusID                        int                `json:"statusId"`
	StatusName                      string             `json:"statusName"`
	StatusMeta                      StatusMeta         `json:"statusMeta"`
	PriorityProcessingTimeMinutes   int                `json:"priorityProcessingTimeMinutes"`
	PriorityFadeTimeMinutes         int                `json:"priorityFadeTimeMinutes"`
	PriorityWeight                  int                `json:"priorityWeight"`
	TotalActiveTime                 int                `json:"totalActiveTime"`
	StatusIsStart                   bool               `json:"statusIsStart"`
	StatusTrackTime                 bool               `json:"statusTrackTime"`
	StatusIsFinal                   bool               `json:"statusIsFinal"`
	StatusChangedAt                 string             `json:"statusChangedAt"`
	IsPainted                       bool               `json:"isPainted"`
	ManagerGroupID                  int                `json:"managerGroupId"`
	ManagerGroupName                string             `json:"managerGroupName"`
	ManagerGroupColor               string             `json:"managerGroupColor"`
	ResponsibleManagerID            int                `json:"responsibleManagerId"`
	ResponsibleManagerFirstName     string             `json:"responsibleManagerFirstName"`
	ResponsibleManagerLastName      string             `json:"responsibleManagerLastName"`
	PreviewData                     PreviewData        `json:"previewData"`
	AccountPositionID               int                `json:"accountPositionId"`
	AccountPositionName             string             `json:"accountPositionName"`
	AccountAvatarMinID              int                `json:"accountAvatarMinId"`
	AccountAvatarMinLink            string             `json:"accountAvatarMinLink"`
	ResponsibleManagerPositionID    int                `json:"responsibleManagerPositionId"`
	ResponsibleManagerPositionName  string             `json:"responsibleManagerPositionName"`
	ResponsibleManagerAvatarMinID   int                `json:"responsibleManagerAvatarMinId"`
	ResponsibleManagerAvatarMinLink string             `json:"responsibleManagerAvatarMinLink"`
	FormCellAnswers                 [][]FormCellAnswer `json:"formCellAnswers"`
	StatusCode                      int                `json:"-"` // Не из API, для внутреннего использования
	Message                         string             `json:"-"` // Не из API, для внутреннего использования
}
