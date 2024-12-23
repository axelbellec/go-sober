import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { DrinkHistoryView } from "@/components/views/drink-history-view";
import { Separator } from "@/components/ui/separator";

export default function DrinkLogPage() {
  return (
    <div className="container flex h-screen w-full flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            Log a Drink 🍺
          </h1>
          <p className="text-sm text-muted-foreground">
            Search and select your drink
          </p>
        </div>
        <DrinkLogForm />
        <Separator className="my-4" />
        <DrinkHistoryView />
      </div>
    </div>
  );
}
