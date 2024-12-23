import { BACTimelineView } from "@/components/views/bac-timeline-view";
import { PageLayout } from "@/components/layouts/page-layout";

export default function AnalyticsPage() {
  return (
    <PageLayout
      heading="Blood Alcohol Content Timeline 📊"
      subheading="Track your BAC over time"
    >
      <BACTimelineView />
    </PageLayout>
  );
}
