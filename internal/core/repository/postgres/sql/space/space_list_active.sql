SELECT
    space_id,
    code,
    name,
    kind,
    capacity,
    photo_url,
    color,
    open_time::text AS open_time,
    close_time::text AS close_time,
    slot_step_min,
    is_active
FROM spaces
WHERE is_active
ORDER BY space_id;
