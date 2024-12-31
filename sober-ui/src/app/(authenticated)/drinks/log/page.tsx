import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { DrinkHistoryView } from "@/components/views/drink-history-view";
import { Separator } from "@/components/ui/separator";
import { PageLayout } from "@/components/layouts/page-layout";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { BarChart3 } from "lucide-react";

export default function DrinkLogPage() {
  return (
    <PageLayout
      heading="Track a Drink 🍺"
      subheading="Keep track of what you're drinking"
      className="sm:w-[350px]"
    >
      <DrinkLogForm />
      <Separator className="my-4" />
      <div className="flex-1 overflow-y-auto">
        <DrinkHistoryView />
      </div>
      <Separator className="my-4" />
      <Link href="/analytics" className="block">
        <Button variant="outline" className="w-full">
          <BarChart3 className="mr-2 h-4 w-4" />
          View Analytics
        </Button>
      </Link>
    </PageLayout>
  );
}
