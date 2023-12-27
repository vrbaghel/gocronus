package types

type NotificationActions struct {
	Actions []string `json:"actions"`
}

type Notification struct {
	ID          int    `json:"id" binding:"required"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Action      string `json:"action" binding:"required"`
	Timezone    string `json:"timezone"`
	ScheduledOn string `json:"scheduled_on"`
	Device      string `json:"device"`
	Status      string `json:"status" binding:"required"`
}

type GetNotificationsResponsePayload struct {
	Notifications []Notification `json:"notifications"`
	Pagination    PaginationData `json:"pagination"`
}

type SendNotificationRequestPayload struct {
	Action       string                                    `json:"action" binding:"required"`
	Timezone     string                                    `json:"timezone,omitempty"`
	ScheduledFor string                                    `json:"scheduled_for,omitempty"`
	Device       string                                    `json:"device" binding:"required"`
	ClickAction  string                                    `json:"click_action,omitempty"`
	Category     *SendNotificationRequestCategoryPayload   `json:"category" binding:"required"`
	Navigation   *SendNotificationRequestNavigationPayload `json:"navigation" binding:"required"`
}

type SendNotificationRequestCategoryPayload struct {
	Type string                                      `json:"type" binding:"required"`
	Data *SendNotificationRequestCategoryDataPayload `json:"data" binding:"required"`
}

type SendNotificationRequestNavigationPayload struct {
	Type string                                 `json:"type" binding:"required"`
	Data *SendNotificationRequestNavDataPayload `json:"data" binding:"required"`
}

type SendNotificationRequestCategoryDataPayload struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	ImageURLs []string `json:"img_urls,omitempty"`
	GifURLs   []string `json:"gif_urls,omitempty"`
}

type SendNotificationRequestNavDataPayload struct {
	PackageID   string `json:"package_id,omitempty"`
	PackageName string `json:"package_name,omitempty"`
	OrderID     string `json:"order_id,omitempty"`
	FilterID    string `json:"filter_id,omitempty"`
	ToolID      string `json:"tool_id,omitempty"`
}

type APIError struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type PaginationData struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

type RequestNotificationPayload struct {
	To             string                                `json:"to" binding:"required"`
	MutableContent bool                                  `json:"mutable_content"`
	Notification   *RequestNotificationAdditionalPayload `json:"notification,omitempty"`
	Data           *RequestNotificationDataPayload       `json:"data" binding:"required"`
}

type RequestNotificationAdditionalPayload struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ClickAction string `json:"click_action"`
}

type RequestNotificationDataPayload struct {
	Id          int    `json:"id" binding:"required"`
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
	Source      int    `json:"source"`
	Category    int    `json:"category"`
	NavType     int    `json:"navType"`
	ImageUrls   string `json:"imageUrls,omitempty"`
	GifUrls     string `json:"gifUrls,omitempty"`
	PackageId   string `json:"packageid,omitempty"`
	PackageName string `json:"packageName,omitempty"`
	OrderId     string `json:"orderid,omitempty"`
	FilterId    string `json:"filterid,omitempty"`
	ToolId      string `json:"toolid,omitempty"`
}
