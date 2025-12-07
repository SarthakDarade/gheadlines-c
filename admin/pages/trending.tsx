import DashboardLayout from '@/layouts/DashboardLayout'
import TrendingModule from '@/modules/trending'

export default function TrendingPage() {
    return (
        <DashboardLayout title="Trending News">
            <div className="mb-8">
                <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-orange-500">
                    Trending Stories
                </h1>
                <p className="text-gray-500 mt-1">Curate what's hot right now.</p>
            </div>
            <TrendingModule />
        </DashboardLayout>
    )
}
