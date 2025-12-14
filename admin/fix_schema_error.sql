-- ===================================================
-- FIX: RELOAD SCHEMA CACHE & VERIFY TABLE
-- ===================================================

-- 1. Force PostgREST to refresh its schema cache
-- This fixes the "Could not find the 'name' column... in the schema cache" error
NOTIFY pgrst, 'reload config';

-- 2. Verify Table Columns (Safety Check)
DO $$
BEGIN
    -- Ensure table exists
    CREATE TABLE IF NOT EXISTS public.sources (
        id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- Ensure 'name' exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'name') THEN
        ALTER TABLE public.sources ADD COLUMN name text;
    END IF;

    -- Ensure 'category' exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'category') THEN
        ALTER TABLE public.sources ADD COLUMN category text DEFAULT 'General';
    END IF;

    -- Ensure 'url' exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'url') THEN
        ALTER TABLE public.sources ADD COLUMN url text;
    END IF;
    
    -- Ensure 'logo_url' is GONE
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'sources' AND column_name = 'logo_url') THEN
        ALTER TABLE public.sources DROP COLUMN logo_url;
    END IF;
END $$;
