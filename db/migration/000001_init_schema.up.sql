-- Create accounts table
CREATE TABLE accounts (
    id bigserial PRIMARY KEY,
    owner VARCHAR NOT NULL CHECK (owner <> ''),  -- Disallow empty strings
    balance BIGINT NOT NULL,
    currency VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create entries table
CREATE TABLE entries (
    id bigserial PRIMARY KEY,
    account_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add comment to entries.amount column
COMMENT ON COLUMN entries.amount IS 'can be negative or positive';

-- Create transfers table
CREATE TABLE transfers (
    id bigserial PRIMARY KEY,
    from_account_id BIGINT NOT NULL,
    to_account_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add comment to transfers.amount column
COMMENT ON COLUMN transfers.amount IS 'must be positive';

-- Create indexes
CREATE INDEX account_owner_idx ON accounts(owner);
CREATE INDEX entries_account_id_idx ON entries(account_id);
CREATE INDEX transfers_from_account_id_idx ON transfers(from_account_id);
CREATE INDEX transfers_to_account_id_idx ON transfers(to_account_id);
CREATE INDEX transfers_from_account_id_to_account_id_idx ON transfers(from_account_id, to_account_id);

-- Add foreign key constraints
ALTER TABLE entries ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE transfers ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY (from_account_id) REFERENCES accounts(id);
ALTER TABLE transfers ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY (to_account_id) REFERENCES accounts(id);