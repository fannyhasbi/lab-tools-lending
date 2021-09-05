package types

type (
	ChatSessionStatusType string
	TopicType             string

	ChatSession struct {
		ID        int64                 `json:"id"`
		Status    ChatSessionStatusType `json:"status"`
		UserID    int64                 `json:"user_id"`
		CreatedAt string                `json:"created_at"`
		UpdatedAt string                `json:"updated_at"`
	}

	ChatSessionDetail struct {
		ID            int64     `json:"id"`
		Topic         TopicType `json:"topic"`
		ChatSessionID int64     `json:"chat_session_id"`
		Data          string    `json:"data"`
		CreatedAt     string    `json:"created_at"`
	}
)

var (
	ChatSessionStatus map[string]ChatSessionStatusType = map[string]ChatSessionStatusType{
		"progress": "PROGRESS",
		"complete": "COMPLETE",
	}

	Topic map[string]TopicType = map[string]TopicType{
		"register_init":     "RGR_init",
		"register_confirm":  "RGR_confirm",
		"register_complete": "RGR_complete",

		"borrow_init":    "BRW_init",
		"borrow_date":    "BRW_date",
		"borrow_reason":  "BRW_reason",
		"borrow_confirm": "BRW_confirm",

		"tool_returning_init":     "RET_init",
		"tool_returning_confirm":  "RET_confim",
		"tool_returning_complete": "RET_complete",

		// admin stuffs
		"respond_borrow_init":             "RESPOND_brw_init",
		"respond_borrow_complete":         "RESPOND_brw_complete",
		"respond_tool_returning_init":     "RESPOND_ret_init",
		"respond_tool_returning_complete": "RESPOND_ret_complete",

		"manage_add_init":    "MNG_add_init",
		"manage_add_name":    "MNG_add_name",
		"manage_add_brand":   "MNG_add_brand",
		"manage_add_type":    "MNG_add_type",
		"manage_add_weight":  "MNG_add_weight",
		"manage_add_stock":   "MNG_add_stock",
		"manage_add_info":    "MNG_add_info",
		"manage_add_photo":   "MNG_add_photo",
		"manage_add_confirm": "MNG_add_confirm",

		"manage_edit_init":    "MNG_edit_init",
		"manage_edit_name":    "MNG_edit_name",
		"manage_edit_brand":   "MNG_edit_brand",
		"manage_edit_type":    "MNG_edit_type",
		"manage_edit_weight":  "MNG_edit_weight",
		"manage_edit_stock":   "MNG_edit_stock",
		"manage_edit_info":    "MNG_edit_info",
		"manage_edit_confirm": "MNG_edit_confirm",
	}
)
