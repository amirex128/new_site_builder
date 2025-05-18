-- Create multiple databases
CREATE DATABASE IF NOT EXISTS new_site_builder;
CREATE DATABASE IF NOT EXISTS logging_db;

-- Grant all privileges on these databases to the user
GRANT ALL PRIVILEGES ON new_site_builder.* TO 'amirex128'@'%';
GRANT ALL PRIVILEGES ON logging_db.* TO 'amirex128'@'%';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;
