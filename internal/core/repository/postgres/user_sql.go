package postgres

import _ "embed"

var (
	//go:embed sql/user/user_upsert.sql
	qUserUpsert string
	//go:embed sql/user/user_get_by_id.sql
	qUserGetByID string
	//go:embed sql/user/user_get_by_vk_user_id.sql
	qUserGetByVKUserID string
	//go:embed sql/user/user_update_notifications.sql
	qUserUpdateNotifications string
)
