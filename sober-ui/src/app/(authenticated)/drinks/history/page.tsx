import { DrinkHistoryView } from "@/components/views/drink-history-view";

export default function DrinkHistoryPage() {
  return (
    <div className="container py-6">
      <h1 className="mb-6 text-2xl font-semibold">Drink History</h1>
      <DrinkHistoryView />
    </div>
  );
}
