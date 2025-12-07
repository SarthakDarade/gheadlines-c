import DashboardLayout from '@/layouts/DashboardLayout'
import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { Plus, Search, Edit, Trash2, Eye, Filter } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { Toaster, toast } from 'react-hot-toast'

export default function Articles() {
    const [articles, setArticles] = useState<any[]>([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')
    const [filter, setFilter] = useState('all') // 'all', 'latest', 'last_60'
    const router = useRouter()

    useEffect(() => {
        fetchArticles()
    }, [search, filter])

    const fetchArticles = async () => {
        setLoading(true)
        let query = supabase
            .from('articles')
            .select('*', { count: 'exact' })
            .order('created_at', { ascending: false })

        if (search) {
            query = query.ilike('title', `%${search}%`)
        }

        if (filter === 'last_60') {
            const sixtyMinsAgo = new Date(Date.now() - 60 * 60 * 1000).toISOString()
            query = query.gte('created_at', sixtyMinsAgo)
        }

        const { data, error } = await query.limit(50)

        if (error) {
            console.error("Error fetching articles:", error)
            toast.error('Error fetching articles')
        } else {
            setArticles(data || [])
        }
        setLoading(false)
    }

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure you want to delete this article?")) return;

        const { error } = await supabase.from('articles').delete().eq('id', id)
        if (error) {
            toast.error("Failed to delete article")
        } else {
            toast.success("Article deleted")
            setArticles(prev => prev.filter(a => a.id !== id))
        }
    }

    return (
        <DashboardLayout title="Articles">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
                <div>
                    <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">Articles</h1>
                    <p className="text-gray-500 mt-1">Manage all news content</p>
                </div>
                <Link href="/articles/create">
                    <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-lg shadow-blue-600/20 font-medium">
                        <Plus size={18} />
                        Create New Article
                    </button>
                </Link>
            </div>

            <div className="glass-card overflow-hidden border border-white/10">
                <div className="p-4 border-b border-white/10 flex flex-col md:flex-row gap-4 justify-between">
                    <div className="relative flex-1 max-w-sm">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                        <input
                            type="text"
                            placeholder="Search articles..."
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            className="glass-input pl-10"
                        />
                    </div>

                    <div className="flex items-center gap-2">
                        <Filter size={18} className="text-gray-500" />
                        <select
                            value={filter}
                            onChange={(e) => setFilter(e.target.value)}
                            className="glass-input py-2 pr-8"
                        >
                            <option value="all">All Articles</option>
                            <option value="latest">Latest</option>
                            <option value="last_60">Last 60 Minutes</option>
                        </select>
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead className="bg-gray-50/50 dark:bg-white/5 text-left text-sm font-medium text-gray-500 dark:text-gray-400">
                            <tr>
                                <th className="px-6 py-4">Title</th>
                                {/* <th className="px-6 py-4">Author</th> */}
                                {/* Author relationship often flaky if RLS blocks profiles select. Hiding if unused or verify later. Keeping simplistic. */}
                                <th className="px-6 py-4">Status</th>
                                <th className="px-6 py-4">Date</th>
                                <th className="px-6 py-4 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-white/5">
                            {loading ? (
                                <tr>
                                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">Loading articles...</td>
                                </tr>
                            ) : articles.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">No articles found.</td>
                                </tr>
                            ) : (
                                articles.map((article) => (
                                    <tr key={article.id} className="hover:bg-gray-50/50 dark:hover:bg-white/5 transition-colors">
                                        <td className="px-6 py-4">
                                            <div className="font-medium max-w-md truncate" title={article.title}>{article.title}</div>
                                            <div className="text-xs text-gray-500 mt-0.5">{article.slug}</div>
                                        </td>
                                        {/* <td className="px-6 py-4 text-sm text-gray-600 dark:text-gray-300">
                      {article.author?.email || 'Unknown'}
                    </td> */}
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${article.published
                                                ? 'bg-green-100 text-green-800 dark:bg-green-500/20 dark:text-green-400'
                                                : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-500/20 dark:text-yellow-400'
                                                }`}>
                                                {article.published ? 'Published' : 'Draft'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-500">
                                            {new Date(article.created_at).toLocaleString()}
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2">
                                                {/* <button className="p-2 hover:bg-gray-100 dark:hover:bg-white/10 rounded-lg text-gray-500 transition-colors">
                          <Eye size={18} />
                        </button> */}
                                                <Link href={`/articles/${article.id}`}>
                                                    <button className="p-2 hover:bg-blue-50 dark:hover:bg-blue-500/10 rounded-lg text-blue-600 transition-colors">
                                                        <Edit size={18} />
                                                    </button>
                                                </Link>
                                                <button
                                                    onClick={() => handleDelete(article.id)}
                                                    className="p-2 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg text-red-600 transition-colors">
                                                    <Trash2 size={18} />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </DashboardLayout>
    )
}
