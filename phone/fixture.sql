DROP TABLE IF EXISTS phone_numbers;

CREATE TABLE phone_numbers (
    id SERIAL,
    number VARCHAR(40) NOT NULL
);

INSERT INTO phone_numbers (number) VALUES
('1234567890'),
('123 456 7891'),
('(123) 456 7892'),
('(123) 456-7893'),
('123-456-7894'),
('123-456-7890'),
('1234567892'),
('(123)456-7892');