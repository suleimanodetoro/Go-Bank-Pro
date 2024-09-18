DROP INDEX IF EXISTS public.transfers_to_account_id_idx;
DROP INDEX IF EXISTS public.transfers_from_account_id_to_account_id_idx;
DROP INDEX IF EXISTS public.transfers_from_account_id_idx;
DROP INDEX IF EXISTS public.entries_account_id_idx;
DROP INDEX IF EXISTS public.account_owner_idx;

ALTER TABLE ONLY public.transfers DROP CONSTRAINT IF EXISTS transfers_to_account_id_fkey;
ALTER TABLE ONLY public.transfers DROP CONSTRAINT IF EXISTS transfers_from_account_id_fkey;
ALTER TABLE ONLY public.entries DROP CONSTRAINT IF EXISTS entries_account_id_fkey;

DROP TABLE IF EXISTS public.transfers;
DROP SEQUENCE IF EXISTS public.transfers_id_seq;

DROP TABLE IF EXISTS public.entries;
DROP SEQUENCE IF EXISTS public.entries_id_seq;

DROP TABLE IF EXISTS public.accounts;
DROP SEQUENCE IF EXISTS public.account_id_seq;
