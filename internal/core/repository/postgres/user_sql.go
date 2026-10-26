package postgres

import _ "embed"

var (
	//go:embed sql/user/upsert.sql
	qUserUpsert string
	//go:embed sql/user/get_by_id.sql
	qUserGetByID string
	//go:embed sql/user/get_by_vk_user_id.sql
	qUserGetByVKUserID string
	//go:embed sql/user/update_notifications.sql
	qUserUpdateNotifications string
)
