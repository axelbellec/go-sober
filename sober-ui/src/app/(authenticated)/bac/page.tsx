import { AlcoholLevelTimelineView } from "@/components/views/alcohol-level-timeline-view";
import { PageLayout } from "@/components/layouts/page-layout";

export default function AnalyticsPage() {
  return (
    <PageLayout
      heading="Your BAC Timeline 📊"
      subheading="Monitor your blood alcohol content over time"
    >
      <AlcoholLevelTimelineView />
    </PageLayout>
  );
}
