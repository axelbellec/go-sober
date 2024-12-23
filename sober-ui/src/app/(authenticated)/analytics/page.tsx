import { BACTimelineView } from "@/components/views/bac-timeline-view";

export default function AnalyticsPage() {
  return (
    <div className="container flex h-screen w-full flex-col">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 py-6">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            Blood Alcohol Content Timeline 📊
          </h1>
          <p className="text-sm text-muted-foreground">
            Track your BAC over time
          </p>
        </div>
        <BACTimelineView />
      </div>
    </div>
  );
}
