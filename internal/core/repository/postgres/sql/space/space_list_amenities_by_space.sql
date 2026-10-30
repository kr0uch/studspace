SELECT
    sa.space_id,
    a.amenity_id,
    a.name,
    a.icon
FROM space_amenities sa
JOIN amenities a ON a.amenity_id = sa.amenity_id
WHERE sa.space_id = $1
ORDER BY a.amenity_id;
