-- GRANT ALL ACCESS TO ADMINS
-- Explicitly ensure admins can do everything on all tables

-- 1. Profiles (Users)
DROP POLICY IF EXISTS "Admins can do everything on profiles" ON public.profiles;
CREATE POLICY "Admins can do everything on profiles"
    ON public.profiles FOR ALL
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- 2. Articles
DROP POLICY IF EXISTS "Admins can do everything on articles" ON public.articles;
CREATE POLICY "Admins can do everything on articles"
    ON public.articles FOR ALL
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- 3. Categories
DROP POLICY IF EXISTS "Admins can do everything on categories" ON public.categories;
CREATE POLICY "Admins can do everything on categories"
    ON public.categories FOR ALL
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- 4. Trending News
DROP POLICY IF EXISTS "Admins can do everything on trending_news" ON public.trending_news;
CREATE POLICY "Admins can do everything on trending_news"
    ON public.trending_news FOR ALL
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- 5. Live Updates
DROP POLICY IF EXISTS "Admins can do everything on live_updates" ON public.live_updates;
CREATE POLICY "Admins can do everything on live_updates"
    ON public.live_updates FOR ALL
    USING (auth.uid() IN (SELECT id FROM public.admins));

-- 6. Admins (Self Management)
-- Be careful not to break the login, using non-recursive check derived from session usually better, 
-- but since we fixed login logic in code/earlier policy, this is for CRUD.
DROP POLICY IF EXISTS "Superadmins can manage admins" ON public.admins;
CREATE POLICY "Superadmins can manage admins"
    ON public.admins FOR ALL
    USING (
       auth.uid() IN (SELECT id FROM public.admins WHERE role = 'superadmin')
    );
