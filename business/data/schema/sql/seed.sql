INSERT INTO events (id, title, location, event_date)
VALUES ('2944a9cb-ef2d-4632-ac1d-af2b2629d0f2', 'Dinner', 'My House', 'Tonight'),
       ('f88f1860-9a5d-423e-820f-9acb4db3030e', 'Go Programming Lesson', 'At School', 'Tomorrow'),
       ('4cb393fb-dd19-469e-a52c-22a12c0a98df', 'Company Picnic', 'At the Park', 'Saturday')
ON CONFLICT DO NOTHING
