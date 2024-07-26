INSERT INTO users (username, email, "password")
    SELECT 'test',
            'test@test.com',
            'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3' /* 123 */
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username='test'
);