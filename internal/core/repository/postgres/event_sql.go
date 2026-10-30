package postgres

import _ "embed"

var (
	//go:embed sql/event/event_get_view.sql
	qEventGetView string
	//go:embed sql/event/event_list_views.sql
	qEventListViews string
	//go:embed sql/event/event_list_in_range.sql
	qEventListInRange string
	//go:embed sql/event/event_list_overlapping.sql
	qEventListOverlapping string
	//go:embed sql/event/event_busy_days.sql
	qEventBusyDays string
	//go:embed sql/event/event_create.sql
	qEventCreate string
	//go:embed sql/event/event_update.sql
	qEventUpdate string
	//go:embed sql/event/event_cancel.sql
	qEventCancel string
	//go:embed sql/event/event_upsert_registration.sql
	qEventUpsertRegistration string
	//go:embed sql/event/event_cancel_registration.sql
	qEventCancelRegistration string
	//go:embed sql/event/event_add_waitlist.sql
	qEventAddWaitlist string
	//go:embed sql/event/event_remove_waitlist.sql
	qEventRemoveWaitlist string
	//go:embed sql/event/event_claim_waitlist.sql
	qEventClaimWaitlist string
)
