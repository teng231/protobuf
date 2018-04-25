CREATE TABLE starfleet (
	name TEXT NOT NULL,
	passengers INT(6),
	mission TEXT,
	departure_time_of_ship TIMESTAMP
);

INSERT INTO starfleet (name, passengers, mission, departure_time_of_ship)
VALUES ("USS Enterprise", 654, NULL, now());

INSERT INTO starfleet (name, passengers, mission, departure_time_of_ship)
VALUES ("USS Discovery", NULL, "spore drive development", NULL);
