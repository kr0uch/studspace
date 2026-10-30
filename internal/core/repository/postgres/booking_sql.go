package postgres

import _ "embed"

var (
	//go:embed sql/booking/booking_create.sql
	qBookingCreate string
	//go:embed sql/booking/booking_get_by_id.sql
	qBookingGetByID string
	//go:embed sql/booking/booking_list_by_user_active.sql
	qBookingListByUserActive string
	//go:embed sql/booking/booking_list_by_user_history.sql
	qBookingListByUserHistory string
	//go:embed sql/booking/booking_list_admin.sql
	qBookingListAdmin string
	//go:embed sql/booking/booking_list_admin_approved.sql
	qBookingListAdminApproved string
	//go:embed sql/booking/booking_count_pending_by_user.sql
	qBookingCountPendingByUser string
	//go:embed sql/booking/booking_count_by_status.sql
	qBookingCountByStatus string
	//go:embed sql/booking/booking_list_overlapping.sql
	qBookingListOverlapping string
	//go:embed sql/booking/booking_list_in_range.sql
	qBookingListInRange string
	//go:embed sql/booking/booking_busy_days.sql
	qBookingBusyDays string
	//go:embed sql/booking/booking_insert_status_log.sql
	qBookingInsertStatusLog string
	//go:embed sql/booking/booking_approve.sql
	qBookingApprove string
	//go:embed sql/booking/booking_reject.sql
	qBookingReject string
	//go:embed sql/booking/booking_withdraw.sql
	qBookingWithdraw string
	//go:embed sql/booking/booking_cancel.sql
	qBookingCancel string
	//go:embed sql/booking/booking_return_to_review.sql
	qBookingReturnToReview string
	//go:embed sql/booking/booking_reschedule.sql
	qBookingReschedule string
)
