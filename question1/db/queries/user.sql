
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id, name, phone_number, otp, otp_expiration_time;

SELECT * FROM users WHERE phone_number = $1;

UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE phone_number = $3 RETURNING id, name, phone_number, otp, otp_expiration_time;

SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1;