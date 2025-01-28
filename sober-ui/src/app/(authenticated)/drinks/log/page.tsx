import { ConsumptionHistoryView } from "@/components/views/consumption-history-view";
import { Separator } from "@/components/ui/separator";
import { PageLayout } from "@/components/layouts/page-layout";

export default function DrinkLogPage() {
  return (
    <PageLayout
      heading="Track a Drink 🍺"
      subheading="Keep track of what you're drinking"
      className="sm:w-[350px]"
    >
      <Separator className="my-4" />
      <div className="flex-1 overflow-y-auto">
        <ConsumptionHistoryView />
      </div>
    </PageLayout>
  );
}
