import { useRouter } from 'next/router'
import Link from 'next/link'
import { motion } from 'framer-motion'
import {
    LayoutDashboard,
    FileText,
    Radio,
    TrendingUp,
    Users,
    Briefcase,
    Mail,
    Settings,
    LogOut,
    ChevronRight,
    Zap,
    Globe
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'

const menuItems = [
    { path: '/', label: 'Overview', icon: LayoutDashboard },
    { path: '/articles', label: 'Content', icon: FileText },
    { path: '/sources', label: 'Sources', icon: Globe },
    { path: '/live', label: 'Live News', icon: Radio },
    { path: '/breaking', label: 'Breaking Alerts', icon: Zap },
    { path: '/trending', label: 'Trending', icon: TrendingUp },
    { path: '/users', label: 'Users', icon: Users },
    { path: '/careers', label: 'Careers', icon: Briefcase },
    { path: '/communications', label: 'Communications', icon: Mail },
    { path: '/settings', label: 'Settings', icon: Settings },
]

export default function Sidebar() {
    const router = useRouter()
    const { signOut } = useAuth()

    return (
        <aside className="fixed left-0 top-0 h-screen w-64 glass border-r border-white/10 flex flex-col z-50 transition-all duration-300">
            <div className="h-16 flex items-center px-6 border-b border-white/10">
                <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center mr-3">
                    <span className="text-white font-bold text-lg">T</span>
                </div>
                <span className="font-bold text-lg tracking-tight">TroyGH Admin</span>
            </div>

            <nav className="flex-1 overflow-y-auto py-6 px-3 space-y-1">
                {menuItems.map((item) => {
                    const isActive = router.pathname === item.path || (item.path !== '/' && router.pathname.startsWith(item.path))
                    const Icon = item.icon

                    return (
                        <Link key={item.path} href={item.path}>
                            <div
                                className={`flex items-center justify-between px-3 py-2.5 rounded-lg transition-all duration-200 group cursor-pointer ${isActive
                                    ? 'bg-blue-600/10 text-blue-600 dark:text-blue-400'
                                    : 'text-gray-600 dark:text-gray-400 hover:bg-white/50 dark:hover:bg-white/5 hover:text-gray-900 dark:hover:text-gray-200'
                                    }`}
                            >
                                <div className="flex items-center gap-3">
                                    <Icon size={18} strokeWidth={isActive ? 2.5 : 2} />
                                    <span className={`text-sm font-medium ${isActive ? 'font-semibold' : ''}`}>{item.label}</span>
                                </div>
                                {isActive && (
                                    <motion.div layoutId="activeSim" className="w-1.5 h-1.5 rounded-full bg-blue-600" />
                                )}
                            </div>
                        </Link>
                    )
                })}
            </nav>

            <div className="p-4 border-t border-white/10 space-y-2">
                <button
                    onClick={signOut}
                    className="w-full flex items-center gap-3 px-3 py-2.5 text-gray-600 dark:text-gray-400 hover:bg-red-500/10 hover:text-red-600 rounded-lg transition-colors text-sm font-medium"
                >
                    <LogOut size={18} />
                    <span>Sign Out</span>
                </button>
            </div>
        </aside>
    )
}
