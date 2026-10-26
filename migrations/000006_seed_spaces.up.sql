INSERT INTO spaces (code, name, kind, capacity, color) VALUES
    ('A', 'Пространство A', 'Коворкинг', 25, '#6C5CE7'),
    ('B', 'Пространство B', 'Переговорная', 8, '#C5E84A');

INSERT INTO amenities (name, icon) VALUES
    ('Проектор', 'projector'),
    ('Wi-Fi', 'wifi'),
    ('Маркерная доска', 'whiteboard'),
    ('Круглый стол', 'round-table');

INSERT INTO space_amenities (space_id, amenity_id)
SELECT s.space_id, a.amenity_id
FROM spaces s
JOIN amenities a ON (s.code, a.name) IN (
    ('A', 'Проектор'),
    ('A', 'Wi-Fi'),
    ('B', 'Маркерная доска'),
    ('B', 'Круглый стол')
);
