import { useState, useEffect } from 'react'
import { Search, Bell, Sun, Moon, User as UserIcon } from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { motion } from 'framer-motion'

export default function Topbar() {
    const { profile } = useAuth()
    const [isDark, setIsDark] = useState(false)

    useEffect(() => {
        // Check local storage or system
        if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            setIsDark(true)
            document.documentElement.classList.add('dark')
        } else {
            setIsDark(false)
            document.documentElement.classList.remove('dark')
        }
    }, [])

    const toggleTheme = () => {
        if (isDark) {
            document.documentElement.classList.remove('dark')
            localStorage.theme = 'light'
            setIsDark(false)
        } else {
            document.documentElement.classList.add('dark')
            localStorage.theme = 'dark'
            setIsDark(true)
        }
    }

    return (
        <header className="h-16 glass border-b border-white/10 sticky top-0 z-40 flex items-center justify-between px-6 transition-all duration-300">
            {/* Search Bar */}
            <div className="flex-1 max-w-xl">
                <div className="relative group">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-blue-500 transition-colors" size={18} />
                    <input
                        type="text"
                        placeholder="Global Search (CMD + K)"
                        className="w-full bg-white/50 dark:bg-black/20 border border-gray-200 dark:border-gray-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                    />
                    <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1">
                        <span className="text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded text-gray-500">⌘</span>
                        <span className="text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded text-gray-500">K</span>
                    </div>
                </div>
            </div>

            {/* Right Actions */}
            <div className="flex items-center gap-4">
                <button
                    onClick={toggleTheme}
                    className="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-white/5 rounded-full transition-colors"
                >
                    {isDark ? <Sun size={20} /> : <Moon size={20} />}
                </button>

                <div className="relative">
                    <button className="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-white/5 rounded-full transition-colors">
                        <Bell size={20} />
                        <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border border-white dark:border-black"></span>
                    </button>
                </div>

                <div className="h-8 w-[1px] bg-gray-200 dark:bg-gray-700 mx-2"></div>

                <div className="flex items-center gap-3">
                    <div className="text-right hidden sm:block">
                        <p className="text-sm font-medium leading-none">{profile?.email || 'Admin User'}</p>
                        <p className="text-xs text-gray-500 mt-1">{profile?.role || 'Superadmin'}</p>
                    </div>
                    <div className="w-9 h-9 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-medium shadow-md">
                        {profile?.email?.[0].toUpperCase() || 'A'}
                    </div>
                </div>
            </div>
        </header>
    )
}
