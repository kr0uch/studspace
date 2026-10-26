package postgres

import _ "embed"

var (
	//go:embed sql/space/list_active.sql
	qSpaceListActive string
	//go:embed sql/space/get_by_id.sql
	qSpaceGetByID string
	//go:embed sql/space/list_amenities.sql
	qSpaceListAmenities string
	//go:embed sql/space/list_amenities_by_space.sql
	qSpaceListAmenitiesBySpace string
)
