package postgres

import _ "embed"

var (
	//go:embed sql/booking/create.sql
	qBookingCreate string
	//go:embed sql/booking/get_by_id.sql
	qBookingGetByID string
	//go:embed sql/booking/list_by_user_active.sql
	qBookingListByUserActive string
	//go:embed sql/booking/list_by_user_history.sql
	qBookingListByUserHistory string
	//go:embed sql/booking/list_admin.sql
	qBookingListAdmin string
	//go:embed sql/booking/list_admin_approved.sql
	qBookingListAdminApproved string
	//go:embed sql/booking/count_pending_by_user.sql
	qBookingCountPendingByUser string
	//go:embed sql/booking/list_overlapping.sql
	qBookingListOverlapping string
	//go:embed sql/booking/list_in_range.sql
	qBookingListInRange string
	//go:embed sql/booking/busy_days.sql
	qBookingBusyDays string
	//go:embed sql/booking/insert_status_log.sql
	qBookingInsertStatusLog string
	//go:embed sql/booking/approve.sql
	qBookingApprove string
	//go:embed sql/booking/reject.sql
	qBookingReject string
	//go:embed sql/booking/withdraw.sql
	qBookingWithdraw string
	//go:embed sql/booking/cancel.sql
	qBookingCancel string
	//go:embed sql/booking/return_to_review.sql
	qBookingReturnToReview string
	//go:embed sql/booking/reschedule.sql
	qBookingReschedule string
)
