-- ============================================
-- 1. ADMINS TABLE SETUP
-- ============================================

CREATE TABLE IF NOT EXISTS public.admins (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    role TEXT CHECK (role IN ('superadmin', 'editor', 'reporter')) DEFAULT 'reporter',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_login TIMESTAMPTZ,
    ip_address TEXT
);

-- Enable RLS
ALTER TABLE public.admins ENABLE ROW LEVEL SECURITY;

-- Policies
-- Only admins can view the admin table
CREATE POLICY "Admins can view other admins" 
    ON public.admins FOR SELECT 
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- ============================================
-- 2. CREATE FIRST ADMIN USER (Seed)
-- ============================================

-- INSTRUCTIONS:
-- 1. Go to Authentication > Users in Supabase Dashboard.
-- 2. Create a new user with email: 'admin@troygh.com' and a secure password.
-- 3. Copy the 'User UID' of the new user.
-- 4. Run the following SQL, replacing 'YOUR_USER_UID_HERE' with the copied UID.

/*
INSERT INTO public.admins (id, email, role)
VALUES ('YOUR_USER_UID_HERE', 'admin@troygh.com', 'superadmin');
*/

-- ============================================
-- 3. ENSURE ARTICLES TABLE HAS EXTRA FIELDS
-- ============================================

-- Add columns if they don't exist (Idempotent)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'category_id') THEN
        ALTER TABLE public.articles ADD COLUMN category_id UUID REFERENCES public.categories(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'slug') THEN
        ALTER TABLE public.articles ADD COLUMN slug TEXT UNIQUE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'tags') THEN
        ALTER TABLE public.articles ADD COLUMN tags TEXT[];
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'published') THEN
        ALTER TABLE public.articles ADD COLUMN published BOOLEAN DEFAULT false;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'meta_title') THEN
        ALTER TABLE public.articles ADD COLUMN meta_title TEXT;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'meta_description') THEN
        ALTER TABLE public.articles ADD COLUMN meta_description TEXT;
    END IF;
     IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'articles' AND column_name = 'image_url') THEN
        ALTER TABLE public.articles ADD COLUMN image_url TEXT;
    END IF;
END $$;
