import DashboardLayout from '@/layouts/DashboardLayout'
import BreakingNewsModule from '@/modules/breaking'

export default function BreakingNewsPage() {
    return (
        <DashboardLayout title="Breaking News">
            <BreakingNewsModule />
        </DashboardLayout>
    )
}
