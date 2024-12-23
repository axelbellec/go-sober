import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { DrinkHistoryView } from "@/components/views/drink-history-view";
import { Separator } from "@/components/ui/separator";
import { PageLayout } from "@/components/layouts/page-layout";

export default function DrinkLogPage() {
  return (
    <PageLayout
      heading="Log a Drink 🍺"
      subheading="Search and select your drink"
      className="sm:w-[350px]"
    >
      <DrinkLogForm />
      <Separator className="my-4" />
      <div className="flex-1 overflow-y-auto">
        <DrinkHistoryView />
      </div>
    </PageLayout>
  );
}
