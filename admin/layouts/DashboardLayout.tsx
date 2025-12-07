import { ReactNode } from 'react'
import Sidebar from '@/components/Sidebar'
import Topbar from '@/components/Topbar'
import { motion } from 'framer-motion'
import { useAuth } from '@/hooks/useAuth'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import Head from 'next/head'

interface DashboardLayoutProps {
    children: ReactNode
    title?: string
}

export default function DashboardLayout({ children, title }: DashboardLayoutProps) {
    const { user, isLoading } = useAuth()
    const router = useRouter()

    useEffect(() => {
        if (!isLoading && !user) {
            router.push('/login')
        }
    }, [user, isLoading, router])

    if (isLoading || !user) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <div className="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full"></div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background text-foreground flex">
            <Head>
                <title>{title ? `${title} | TroyGH Admin` : 'TroyGH Admin'}</title>
            </Head>

            <Sidebar />

            <div className="flex-1 ml-64 flex flex-col min-h-screen">
                <Topbar />

                <main className="flex-1 p-6 overflow-x-hidden">
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.4 }}
                    >
                        {children}
                    </motion.div>
                </main>
            </div>
        </div>
    )
}
