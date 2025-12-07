import { createClient } from '@supabase/supabase-js'

// Hardcoded for reliability in this specific environment context
const supabaseUrl = 'https://dcqppqeydxdnbscharjn.supabase.co'
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRjcXBwcWV5ZHhkbmJzY2hhcmpuIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQwNzk5NjgsImV4cCI6MjA3OTY1NTk2OH0.Mxt9n-mGcvZb4QIukohuiJA4fkiOxtKgdwPVgLhaNmk'

if (!supabaseUrl || !supabaseAnonKey) {
    console.warn('Supabase URL or Key is missing in environment variables.')
}

export const supabase = createClient(supabaseUrl, supabaseAnonKey)
