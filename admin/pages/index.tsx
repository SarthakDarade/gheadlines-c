import DashboardLayout from '@/layouts/DashboardLayout'
import StatsCard from '@/components/StatsCard'
import { Users, FileText, Activity, Globe, Zap } from 'lucide-react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import { useEffect, useState } from 'react'
import { supabase } from '@/lib/supabaseClient'
import Link from 'next/link'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

export default function Dashboard() {
  const [stats, setStats] = useState({
    users: 0,
    articles: 0,
    activeVisitors: 12, // Mock until Realtime implemented
    todayViews: 0
  })

  const [recentArticles, setRecentArticles] = useState<any[]>([])

  // Basic Chart Configuration
  const chartData = {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datasets: [
      {
        label: 'Traffic',
        data: [1200, 1900, 3000, 5000, 2000, 3000, 4500],
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        fill: true,
        tension: 0.4,
      },
      {
        label: 'New Users',
        data: [200, 400, 300, 800, 400, 600, 700],
        borderColor: 'rgb(168, 85, 247)',
        backgroundColor: 'rgba(168, 85, 247, 0.1)',
        fill: true,
        tension: 0.4,
      }
    ],
  }

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          color: '#94a3b8' // Slate 400
        }
      },
    },
    scales: {
      y: {
        grid: {
          color: 'rgba(255, 255, 255, 0.1)'
        },
        ticks: {
          color: '#94a3b8'
        }
      },
      x: {
        grid: {
          display: false
        },
        ticks: {
          color: '#94a3b8'
        }
      }
    }
  }

  useEffect(() => {
    async function fetchData() {
      // 1. Fetch Stats
      const [
        { count: userCount },
        { count: articleCount },
      ] = await Promise.all([
        supabase.from('profiles').select('*', { count: 'exact', head: true }),
        supabase.from('articles').select('*', { count: 'exact', head: true }),
      ])

      setStats(prev => ({
        ...prev,
        users: userCount || 0,
        articles: articleCount || 0
      }))

      // 2. Fetch Recent Articles for right column
      const { data: recents } = await supabase
        .from('articles')
        .select('id, title, created_at')
        .order('created_at', { ascending: false })
        .limit(6)

      if (recents) setRecentArticles(recents)
    }

    fetchData()

    // Subscribe to realtime changes
    const subscription = supabase
      .channel('public:all')
      .on('postgres_changes', { event: '*', schema: 'public', table: 'profiles' }, () => fetchData())
      .on('postgres_changes', { event: '*', schema: 'public', table: 'articles' }, () => fetchData())
      .subscribe()

    return () => {
      subscription.unsubscribe()
    }
  }, [])

  return (
    <DashboardLayout title="Overview">
      <div className="mb-8">
        <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">
          Dashboard Overview
        </h1>
        <p className="text-gray-500 mt-2">Welcome back, here's what's happening today.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatsCard
          title="Total Users"
          value={stats.users}
          change="+12% from last week"
          changeType="positive"
          icon={Users}
          color="bg-blue-500 text-blue-500" // Passing classes
        />
        <StatsCard
          title="Total Articles"
          value={stats.articles}
          change="+5 new today"
          changeType="neutral"
          icon={FileText}
          color="bg-purple-500 text-purple-500"
        />
        <StatsCard
          title="Active Visitors"
          value={stats.activeVisitors}
          change="Live now"
          changeType="positive"
          icon={Activity}
          color="bg-green-500 text-green-500"
        />
        <StatsCard
          title="Total Traffic"
          value="45.2k"
          change="+3% overflow"
          changeType="positive"
          icon={Globe}
          color="bg-orange-500 text-orange-500"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 glass-card p-6 border border-white/10">
          <h2 className="text-xl font-bold mb-6">Traffic Analytics</h2>
          <div className="h-[300px]">
            <Line options={chartOptions} data={chartData} />
          </div>
        </div>

        <div className="glass-card p-6 border border-white/10">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-bold">Recent Articles</h2>
            <Link href="/articles" className="text-xs text-blue-500 hover:underline">View All</Link>
          </div>

          <div className="space-y-4">
            {recentArticles.length === 0 ? (
              <p className="text-gray-500 text-sm">No recent articles.</p>
            ) : (
              recentArticles.map((article) => (
                <Link key={article.id} href={`/articles/${article.id}`}>
                  <div className="flex items-center gap-4 p-3 hover:bg-white/50 dark:hover:bg-white/5 rounded-lg transition-colors cursor-pointer group">
                    <div className="w-10 h-10 rounded-full bg-blue-500/10 flex items-center justify-center text-blue-500 group-hover:bg-blue-600 group-hover:text-white transition-colors">
                      <Zap size={18} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate text-gray-900 dark:text-gray-100">{article.title}</p>
                      <p className="text-xs text-gray-500">
                        {new Date(article.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                      </p>
                    </div>
                  </div>
                </Link>
              ))
            )}
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
