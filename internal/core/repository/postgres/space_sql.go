package postgres

import _ "embed"

var (
	//go:embed sql/space/space_list_active.sql
	qSpaceListActive string
	//go:embed sql/space/space_get_by_id.sql
	qSpaceGetByID string
	//go:embed sql/space/space_list_amenities.sql
	qSpaceListAmenities string
	//go:embed sql/space/space_list_amenities_by_space.sql
	qSpaceListAmenitiesBySpace string
)
