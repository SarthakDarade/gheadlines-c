import { motion } from 'framer-motion'
import { LucideIcon } from 'lucide-react'

interface StatsCardProps {
    title: string
    value: string | number
    change?: string
    changeType?: 'positive' | 'negative' | 'neutral'
    icon: LucideIcon
    color: string
}

export default function StatsCard({ title, value, change, changeType = 'neutral', icon: Icon, color }: StatsCardProps) {
    return (
        <motion.div
            whileHover={{ y: -5 }}
            className="glass p-6 rounded-2xl border border-white/10 shadow-sm relative overflow-hidden group"
        >
            <div className={`absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity ${color}`}>
                <Icon size={100} />
            </div>

            <div className="flex justify-between items-start mb-4">
                <div className={`p-3 rounded-xl ${color} bg-opacity-10 text-white`}>
                    <Icon size={24} className={color.replace('text-', 'text-opacity-100 ')} />
                    {/* Note: In Tailwind, utility classes like 'text-blue-500' work. If passed 'text-blue-500', we need to handle bg opacity or just use the color class directly on Icon if properly passed. 
              Better to pass specific classes or verify usage. 
              Let's assume 'color' is like 'bg-blue-500' or similar.
              Actually, let's just make 'color' a class string for text/bg.
          */}
                </div>
                {change && (
                    <span className={`text-xs font-medium px-2 py-1 rounded-full ${changeType === 'positive' ? 'bg-green-500/10 text-green-600' :
                            changeType === 'negative' ? 'bg-red-500/10 text-red-600' : 'bg-gray-500/10 text-gray-600'
                        }`}>
                        {change}
                    </span>
                )}
            </div>

            <div>
                <h3 className="text-gray-500 dark:text-gray-400 text-sm font-medium">{title}</h3>
                <p className="text-3xl font-bold mt-1 text-gray-900 dark:text-white">{value}</p>
            </div>
        </motion.div>
    )
}
