import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { supabase } from '@/lib/supabaseClient'
import { TrendingUp, Plus, X, UploadCloud } from 'lucide-react'
import { toast } from 'react-hot-toast'

interface TrendingItem {
    id: string
    title: string
    summary: string
    category: string
    image_url?: string
}

export default function TrendingModule() {
    const [trending, setTrending] = useState<TrendingItem[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        fetchTrending()
    }, [])

    const fetchTrending = async () => {
        const { data } = await supabase.from('trending_news').select('*').order('created_at', { ascending: false })
        if (data) setTrending(data)
        setLoading(false)
    }

    const handleDelete = async (id: string) => {
        setTrending(prev => prev.filter(t => t.id !== id))
        await supabase.from('trending_news').delete().eq('id', id)
        toast.success("Removed from trending")
    }

    // Expanded form state
    const [showForm, setShowForm] = useState(false)
    const [newItem, setNewItem] = useState({
        title: '',
        summary: '',
        category: 'General',
        image_url: ''
    })

    const handleAdd = async (e: React.FormEvent) => {
        e.preventDefault()
        const { data, error } = await supabase.from('trending_news').insert([newItem]).select().single()

        if (error) {
            toast.error("Failed to add")
        } else if (data) {
            setTrending([data, ...trending])
            setNewItem({ title: '', summary: '', category: 'General', image_url: '' })
            setShowForm(false)
            toast.success("Added to trending")
        }
    }

    return (
        <div className="space-y-8">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {/* Add Card */}
                {!showForm ? (
                    <div
                        onClick={() => setShowForm(true)}
                        className="glass-card p-6 border border-dashed border-gray-300 dark:border-gray-700 flex flex-col items-center justify-center text-center cursor-pointer hover:bg-gray-50 dark:hover:bg-white/5 transition-colors group h-auto min-h-[300px]"
                    >
                        <div className="w-12 h-12 rounded-full bg-blue-500/10 flex items-center justify-center text-blue-500 mb-3 group-hover:scale-110 transition-transform">
                            <Plus size={24} />
                        </div>
                        <h3 className="font-bold text-gray-900 dark:text-gray-100">Add Trending Story</h3>
                        <p className="text-sm text-gray-500 mt-1">Create custom trending card</p>
                    </div>
                ) : (
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        className="glass-card p-6 border border-blue-500/30 ring-2 ring-blue-500/10"
                    >
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="font-bold">New Trending</h3>
                            <button onClick={() => setShowForm(false)} className="text-gray-400 hover:text-gray-600"><X size={18} /></button>
                        </div>
                        <form onSubmit={handleAdd} className="space-y-3">
                            <input
                                className="glass-input text-sm px-3 py-2"
                                placeholder="Title"
                                value={newItem.title}
                                onChange={e => setNewItem({ ...newItem, title: e.target.value })}
                                required
                            />
                            <textarea
                                className="glass-input text-sm px-3 py-2"
                                placeholder="Summary"
                                rows={2}
                                value={newItem.summary}
                                onChange={e => setNewItem({ ...newItem, summary: e.target.value })}
                            />
                            <input
                                className="glass-input text-sm px-3 py-2"
                                placeholder="Category (e.g. Tech)"
                                value={newItem.category}
                                onChange={e => setNewItem({ ...newItem, category: e.target.value })}
                            />
                            <input
                                type="url"
                                className="glass-input text-sm px-3 py-2"
                                placeholder="Image URL (http...)"
                                value={newItem.image_url}
                                onChange={e => setNewItem({ ...newItem, image_url: e.target.value })}
                            />
                            <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 rounded-lg text-sm font-medium transition-colors">
                                Publish to Trending
                            </button>
                        </form>
                    </motion.div>
                )}

                {trending.map((item, index) => (
                    <motion.div
                        key={item.id}
                        layout
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        className="glass-card p-0 border border-white/10 overflow-hidden group relative flex flex-col"
                    >
                        <div className="h-40 bg-gray-200 dark:bg-gray-800 relative overflow-hidden">
                            {item.image_url ? (
                                <img src={item.image_url} className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" alt="" />
                            ) : (
                                <div className="w-full h-full flex items-center justify-center text-gray-400">
                                    <TrendingUp size={32} opacity={0.2} />
                                </div>
                            )}
                            <div className="absolute top-2 right-2 bg-black/60 text-white text-[10px] font-bold px-2 py-1 rounded backdrop-blur-md shadow-sm">
                                #{index + 1}
                            </div>
                            <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-t from-black/60 to-transparent opacity-60" />
                        </div>

                        <div className="p-5 flex-1 flex flex-col">
                            <span className="text-xs font-bold text-blue-500 uppercase tracking-wider mb-1 block">{item.category}</span>
                            <h3 className="font-bold text-lg leading-tight group-hover:text-blue-500 transition-colors mb-2">{item.title}</h3>
                            <p className="text-sm text-gray-500 dark:text-gray-400 line-clamp-3">{item.summary}</p>
                        </div>

                        <button
                            onClick={() => handleDelete(item.id)}
                            className="absolute top-2 left-2 p-2 bg-white/90 dark:bg-black/50 text-red-500 rounded-lg opacity-0 group-hover:opacity-100 transition-all hover:bg-red-500 hover:text-white shadow-lg transform -translate-y-2 group-hover:translate-y-0"
                        >
                            <Trash2 size={16} />
                        </button>
                    </motion.div>
                ))}
            </div>
        </div>
    )
}

function Trash2(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M3 6h18" />
            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
            <line x1="10" x2="10" y1="11" y2="17" />
            <line x1="14" x2="14" y1="11" y2="17" />
        </svg>
    )
}
