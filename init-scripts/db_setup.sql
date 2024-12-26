create database if not exists reservations;
use reservations;

create table if not exists reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    national_id VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    num_guests INT NOT NULL
);

-- create index idx_reservation_dates on reservations (start_date, end_date);

CREATE TABLE calendar (
    date DATE NOT NULL PRIMARY KEY,
    is_available BOOLEAN NOT NULL DEFAULT TRUE
);

WITH RECURSIVE date_generator AS (
    SELECT '2024-01-01' AS date
    UNION ALL
    SELECT DATE_ADD(date, INTERVAL 1 DAY)
    FROM date_generator
    WHERE date < '2034-12-31'
)
INSERT INTO calendar (date)
SELECT date FROM date_generator;

-- Insert test data
INSERT INTO reservations (
  email, first_name, last_name, national_id, start_date, end_date, num_guests
) VALUES
('john.doe@example.com', 'John', 'Doe', 'ABC123', '2024-12-25', '2024-12-31', 2),
('jane.smith@example.com', 'Jane', 'Smith', 'XYZ789', '2025-01-01', '2025-01-03', 4),
('michael.brown@example.com', 'Michael', 'Brown', 'LMN456', '2025-01-04', '2025-01-05', 1),
('linda.jones@example.com', 'Linda', 'Jones', 'QRS987', '2025-01-06', '2025-01-31', 3);

UPDATE calendar SET is_available = FALSE
WHERE TRUE
AND DATE between '2024-12-25' AND '2024-12-31'
AND DATE between '2025-01-01' AND '2025-01-03'
AND DATE between '2025-01-04' AND '2025-01-05'
AND DATE between '2025-01-06' AND '2025-01-31'
;
