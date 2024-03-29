CREATE DATABASE customers;

USE customers;

CREATE TABLE customer (
	customer_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	uuid VARCHAR(255) NOT NULL,
	email VARCHAR(50) NOT NULL,
	first_name VARCHAR(50) NOT NULL,
	middle_name VARCHAR(50),
	last_name VARCHAR(50) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (customer_id)
);
