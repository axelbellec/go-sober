import { PageLayout } from "@/components/layouts/page-layout";
import { DrinkingTrendsView } from "@/components/views/drinking-trends-view";
import { SobrietyDashboardView } from "@/components/views/sobriety-dashboard-view";
import { SobrietyStatsView } from "@/components/views/sobriety-stats-view";

export default function AnalyticsPage() {
  return (
    <PageLayout
      heading="Analytics 📈"
      subheading="Monitor your drinking habits over time"
    >
      <DrinkingTrendsView />
      <SobrietyDashboardView />
      <SobrietyStatsView />
    </PageLayout>
  );
}
