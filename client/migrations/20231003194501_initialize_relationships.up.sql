CREATE TABLE IF NOT EXISTS credentials (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128) NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    register_date TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT unique_email_cred UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NOT NULL,
    credential_id INTEGER NOT NULL,

    CONSTRAINT unique_credential_id UNIQUE(credential_id),
    CONSTRAINT fk_credentials_users FOREIGN KEY (credential_id) REFERENCES credentials(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    transaction_id VARCHAR(256) NOT NULL,

    CONSTRAINT fk_user_transactions FOREIGN KEY (user_id) REFERENCES users(id)
);