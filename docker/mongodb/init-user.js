// init-user.js

// MongoDB initialization script
// This script creates a database and user

db = db.getSiblingDB('new_site_builder');

// Create user for the database
db.createUser({
  user: 'amirex128',
  pwd: 'mI6G5jd3qNlJQinBOnA2z5SVEawLn4WV',
  roles: [
    { role: 'readWrite', db: 'new_site_builder' },
    { role: 'dbAdmin', db: 'new_site_builder' }
  ]
});

// Create a test collection
db.createCollection('test');
db.test.insertOne({ name: 'Initial document', created_at: new Date() });

print('Database initialization completed successfully');