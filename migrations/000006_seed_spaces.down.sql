DELETE FROM space_amenities
WHERE space_id IN (SELECT space_id FROM spaces WHERE code IN ('A', 'B'));

DELETE FROM amenities
WHERE name IN ('Проектор', 'Wi-Fi', 'Маркерная доска', 'Круглый стол');

DELETE FROM spaces
WHERE code IN ('A', 'B');
