package postgres

import _ "embed"

var (
	//go:embed sql/event/get_view.sql
	qEventGetView string
	//go:embed sql/event/list_views.sql
	qEventListViews string
	//go:embed sql/event/list_in_range.sql
	qEventListInRange string
	//go:embed sql/event/list_overlapping.sql
	qEventListOverlapping string
	//go:embed sql/event/create.sql
	qEventCreate string
	//go:embed sql/event/cancel.sql
	qEventCancel string
	//go:embed sql/event/upsert_registration.sql
	qEventUpsertRegistration string
	//go:embed sql/event/cancel_registration.sql
	qEventCancelRegistration string
	//go:embed sql/event/add_waitlist.sql
	qEventAddWaitlist string
	//go:embed sql/event/remove_waitlist.sql
	qEventRemoveWaitlist string
	//go:embed sql/event/claim_waitlist.sql
	qEventClaimWaitlist string
)
