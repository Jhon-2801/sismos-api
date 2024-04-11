CREATE DATABASE IF NOT EXISTS sismo_db;
USE sismo_db;

CREATE TABLE IF NOT EXISTS events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id VARCHAR(255) UNIQUE NOT NULL,
    magnitude DECIMAL,
    place TEXT NOT NULL,
    event_time TIMESTAMP,
    url TEXT NOT NULL,
    tsunami BOOLEAN,
    mag_type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    longitude DECIMAL,
    latitude DECIMAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    feature_id INT,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feature_id) REFERENCES events(id)
);