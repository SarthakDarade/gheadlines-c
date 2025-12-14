-- ===================================================
-- FIX: CHANNEL_NAME CONSTRAINT ISSUE
-- ===================================================

-- The error 'null value in column "channel_name" ... violates not-null constraint' 
-- implies that 'channel_name' exists and is required, but we are inserting into 'name'.

DO $$
BEGIN
    -- 1. If 'channel_name' exists, we need to handle it.
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'channel_name') THEN
        
        -- A. If 'name' also exists (which our previous scripts tried to ensure):
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'name') THEN
            -- Sync data: Copy channel_name to name for existing rows
            UPDATE public.sources SET name = channel_name WHERE name IS NULL;
        
        -- B. If 'name' does NOT exist for some reason, rename channel_name to name
        ELSE
            ALTER TABLE public.sources RENAME COLUMN channel_name TO name;
        END IF;

        -- 2. CRITICAL: Remove the NOT NULL constraint from channel_name
        -- This allows our new code (which doesn't know about channel_name) to insert NULLs securely.
        ALTER TABLE public.sources ALTER COLUMN channel_name DROP NOT NULL;
        
    END IF;
END $$;

-- 3. Force Cache Reload
NOTIFY pgrst, 'reload config';
