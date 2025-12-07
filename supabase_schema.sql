-- Supabase Schema Setup for GHeadlines
-- Run this in your Supabase SQL Editor

-- ============================================
-- 1. CREATE PROFILES TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS public.profiles (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    username TEXT,
    full_name TEXT,
    avatar_url TEXT,
    website TEXT,
    location TEXT,
    bio TEXT,
    occupation TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================
-- 2. ENABLE ROW LEVEL SECURITY (RLS)
-- ============================================

ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;

-- ============================================
-- 3. CREATE RLS POLICIES FOR PROFILES
-- ============================================

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Profiles are viewable by everyone" ON public.profiles;
DROP POLICY IF EXISTS "Users can insert own profile" ON public.profiles;
DROP POLICY IF EXISTS "Users can update own profile" ON public.profiles;

-- Policy: Everyone can read all profiles
CREATE POLICY "Profiles are viewable by everyone"
    ON public.profiles FOR SELECT
    USING (true);

-- Policy: Users can insert their own profile
CREATE POLICY "Users can insert own profile"
    ON public.profiles FOR INSERT
    WITH CHECK (auth.uid() = id);

-- Policy: Users can update their own profile
CREATE POLICY "Users can update own profile"
    ON public.profiles FOR UPDATE
    USING (auth.uid() = id);

-- ============================================
-- 4. CREATE FUNCTION TO AUTO-CREATE PROFILE
-- ============================================
-- This function automatically creates a profile when a user signs up

CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO public.profiles (id, full_name, created_at, updated_at)
    VALUES (
        NEW.id,
        COALESCE(NEW.raw_user_meta_data->>'full_name', ''),
        NOW(),
        NOW()
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- ============================================
-- 5. CREATE TRIGGER FOR AUTO-PROFILE CREATION
-- ============================================

DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;

CREATE TRIGGER on_auth_user_created
    AFTER INSERT ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION public.handle_new_user();

-- ============================================
-- 6. ENSURE ARTICLES TABLE HAS PROPER RLS
-- ============================================

-- Enable RLS on articles (if not already enabled)
ALTER TABLE public.articles ENABLE ROW LEVEL SECURITY;

-- Drop existing policy if it exists
DROP POLICY IF EXISTS "Articles are viewable by everyone" ON public.articles;

-- Policy: Everyone can read articles
CREATE POLICY "Articles are viewable by everyone"
    ON public.articles FOR SELECT
    USING (true);

-- ============================================
-- 7. CREATE CATEGORIES TABLE (if not exists)
-- ============================================

CREATE TABLE IF NOT EXISTS public.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    color TEXT DEFAULT '#3B82F6',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS on categories
ALTER TABLE public.categories ENABLE ROW LEVEL SECURITY;

-- Policy: Everyone can read categories
DROP POLICY IF EXISTS "Categories are viewable by everyone" ON public.categories;
CREATE POLICY "Categories are viewable by everyone"
    ON public.categories FOR SELECT
    USING (true);

-- ============================================
-- 8. INSERT DEFAULT CATEGORIES
-- ============================================

INSERT INTO public.categories (name, slug, color) VALUES
    ('Technology', 'technology', '#3B82F6'),
    ('Business', 'business', '#F59E0B'),
    ('Politics', 'politics', '#DC2626'),
    ('Sports', 'sports', '#8B5CF6'),
    ('Entertainment', 'entertainment', '#EC4899'),
    ('Science', 'science', '#7C3AED'),
    ('Health', 'health', '#EF4444'),
    ('Environment', 'environment', '#10B981'),
    ('Culture', 'culture', '#EC4899')
ON CONFLICT (slug) DO NOTHING;

-- ============================================
-- 9. CREATE UPDATED_AT TRIGGER FOR PROFILES
-- ============================================

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_profiles_updated_at ON public.profiles;

CREATE TRIGGER update_profiles_updated_at
    BEFORE UPDATE ON public.profiles
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

-- ============================================
-- 10. GRANT PERMISSIONS
-- ============================================

-- Grant usage on schema
GRANT USAGE ON SCHEMA public TO anon, authenticated;

-- Grant access to tables
GRANT SELECT ON public.profiles TO anon, authenticated;
GRANT INSERT, UPDATE ON public.profiles TO authenticated;
GRANT SELECT ON public.articles TO anon, authenticated;
GRANT SELECT ON public.categories TO anon, authenticated;

-- ============================================
-- VERIFICATION QUERIES
-- ============================================

-- Check if profiles table exists
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' AND table_name = 'profiles';

-- Check RLS policies on profiles
SELECT schemaname, tablename, policyname, permissive, roles, cmd, qual 
FROM pg_policies 
WHERE tablename = 'profiles';

-- Check if trigger exists
SELECT trigger_name, event_manipulation, event_object_table 
FROM information_schema.triggers 
WHERE trigger_name = 'on_auth_user_created';

-- ============================================
-- NOTES
-- ============================================

/*
1. This schema supports the GHeadlines application with user profiles
2. RLS is enabled to ensure data security
3. Auto-profile creation trigger ensures every new user gets a profile
4. The application will also create profiles via the CreateProfile method as a fallback
5. Make sure to update your Supabase environment variables in the .env file
*/
