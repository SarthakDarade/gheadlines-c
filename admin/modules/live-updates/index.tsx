import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { supabase } from '@/lib/supabaseClient'
import { Radio, Plus, Trash2, Clock, PlayCircle, Loader2 } from 'lucide-react'
// If you have `react-hot-toast` installed, import it. 
// import { toast } from 'react-hot-toast'

interface LiveUpdate {
    id: string
    headline: string
    source: string
    url: string
    created_at: string
}

export default function LiveUpdatesModule() {
    const [updates, setUpdates] = useState<LiveUpdate[]>([])
    const [loading, setLoading] = useState(true)
    const [headline, setHeadline] = useState('')
    const [source, setSource] = useState('GHeadlines')
    const [url, setUrl] = useState('')

    useEffect(() => {
        fetchUpdates()

        const subscription = supabase
            .channel('live_updates_changes')
            .on('postgres_changes', { event: '*', schema: 'public', table: 'live_updates' }, (payload) => {
                fetchUpdates()
            })
            .subscribe()

        return () => {
            subscription.unsubscribe()
        }
    }, [])

    const fetchUpdates = async () => {
        const { data } = await supabase
            .from('live_updates')
            .select('*')
            .order('created_at', { ascending: false })
            .limit(20) // Keep it snappy

        if (data) setUpdates(data)
        setLoading(false)
    }

    const handleAdd = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!headline) return

        const { error } = await supabase.from('live_updates').insert([{
            headline,
            source,
            url: url || null
        }])

        if (error) {
            alert('Error adding update')
        } else {
            setHeadline('')
            setUrl('')
        }
    }

    const handleDelete = async (id: string) => {
        await supabase.from('live_updates').delete().eq('id', id)
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Control Panel */}
            <div className="lg:col-span-1 space-y-6">
                <div className="glass p-6 rounded-2xl border border-white/10 sticky top-24">
                    <div className="flex items-center gap-3 mb-6">
                        <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center text-red-500 animate-pulse">
                            <Radio size={20} />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold">New Update</h2>
                            <p className="text-xs text-gray-500">Push to homepage instantly</p>
                        </div>
                    </div>

                    <form onSubmit={handleAdd} className="space-y-4">
                        <div>
                            <label className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Headline</label>
                            <textarea
                                value={headline}
                                onChange={(e) => setHeadline(e.target.value)}
                                className="w-full mt-2 p-3 bg-white/50 dark:bg-black/20 border border-gray-200 dark:border-gray-700 rounded-xl focus:ring-2 focus:ring-red-500 outline-none text-sm font-medium"
                                rows={3}
                                placeholder="Breaking news..."
                                required
                            />
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Source</label>
                                <input
                                    type="text"
                                    value={source}
                                    onChange={(e) => setSource(e.target.value)}
                                    className="w-full mt-2 p-3 bg-white/50 dark:bg-black/20 border border-gray-200 dark:border-gray-700 rounded-xl focus:ring-2 focus:ring-red-500 outline-none text-sm"
                                    placeholder="CNN, BBC..."
                                />
                            </div>
                            <div>
                                <label className="text-xs font-semibold text-gray-500 uppercase tracking-wider">URL (Opt)</label>
                                <input
                                    type="url"
                                    value={url}
                                    onChange={(e) => setUrl(e.target.value)}
                                    className="w-full mt-2 p-3 bg-white/50 dark:bg-black/20 border border-gray-200 dark:border-gray-700 rounded-xl focus:ring-2 focus:ring-red-500 outline-none text-sm"
                                    placeholder="https://..."
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            className="w-full py-3 bg-red-600 hover:bg-red-700 text-white rounded-xl font-bold flex items-center justify-center gap-2 transition-all shadow-lg shadow-red-600/20"
                        >
                            <PlayCircle size={18} />
                            Push Live
                        </button>
                    </form>
                </div>
            </div>

            {/* Live Feed */}
            <div className="lg:col-span-2">
                <div className="flex items-center justify-between mb-6">
                    <h2 className="text-xl font-bold flex items-center gap-2">
                        <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />
                        Live Feed
                    </h2>
                    <span className="text-sm text-gray-500 font-mono">
                        {updates.length} active updates
                    </span>
                </div>

                <div className="space-y-4">
                    <AnimatePresence>
                        {loading ? (
                            <div className="flex flex-col items-center justify-center py-20 text-gray-400">
                                <Loader2 className="animate-spin mb-2" size={32} />
                                <p>Connecting to satellite...</p>
                            </div>
                        ) : updates.map((update) => (
                            <motion.div
                                key={update.id}
                                initial={{ opacity: 0, y: -20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, scale: 0.95 }}
                                className="glass p-4 rounded-xl border border-white/5 flex gap-4 group hover:bg-white/5 transition-colors"
                            >
                                <div className="flex flex-col items-center gap-1 pt-1">
                                    <div className="text-[10px] font-bold text-gray-400 font-mono w-16 text-center">
                                        {new Date(update.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                    </div>
                                    <div className="w-[1px] h-full bg-red-500/20 group-hover:bg-red-500/50 transition-colors" />
                                </div>

                                <div className="flex-1">
                                    <h3 className="font-semibold text-gray-900 dark:text-gray-100 leading-snug">
                                        {update.headline}
                                    </h3>
                                    <div className="flex items-center gap-3 mt-2">
                                        <span className="text-xs bg-gray-100 dark:bg-white/10 px-2 py-0.5 rounded text-gray-500">
                                            {update.source}
                                        </span>
                                        {update.url && (
                                            <a href={update.url} target="_blank" rel="noreferrer" className="text-xs text-blue-500 hover:underline">
                                                Source Link
                                            </a>
                                        )}
                                    </div>
                                </div>

                                <button
                                    onClick={() => handleDelete(update.id)}
                                    className="opacity-0 group-hover:opacity-100 p-2 text-gray-400 hover:text-red-500 hover:bg-red-500/10 rounded-lg transition-all self-center"
                                >
                                    <Trash2 size={18} />
                                </button>
                            </motion.div>
                        ))}
                    </AnimatePresence>
                </div>
            </div>
        </div>
    )
}
