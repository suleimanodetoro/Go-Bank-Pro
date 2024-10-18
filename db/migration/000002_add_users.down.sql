-- Drop constraints added during the migration
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Drop the "users" table that was created incorrectly with the extra "id" field
DROP TABLE IF EXISTS "users";
