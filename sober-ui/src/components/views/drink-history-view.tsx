"use client";

import { useEffect, useState } from "react";
import { format, formatDistanceToNow } from "date-fns";
import { DrinkLog } from "@/lib/types/api";
import { fetchWithAuth } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";

export function DrinkHistoryView() {
  const [drinkLogs, setDrinkLogs] = useState<DrinkLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadDrinkLogs();
  }, []);

  async function loadDrinkLogs() {
    try {
      const response = await fetchWithAuth(
        "http://localhost:3000/api/v1/drink-logs"
      );
      if (!response.ok) throw new Error("Failed to load drink history");
      const data = await response.json();
      const sortedDrinkLogs = data.drink_logs.sort(
        (a: DrinkLog, b: DrinkLog) =>
          new Date(b.logged_at).getTime() - new Date(a.logged_at).getTime()
      );
      setDrinkLogs(sortedDrinkLogs);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  }

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-destructive mb-4">{error}</p>
        <Button onClick={() => loadDrinkLogs()}>Retry</Button>
      </div>
    );
  }

  // Group drinks by date
  const groupedDrinks = drinkLogs.reduce((groups, drink) => {
    const date = format(new Date(drink.logged_at), "yyyy-MM-dd");
    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(drink);
    return groups;
  }, {} as Record<string, DrinkLog[]>);

  return (
    <div className="space-y-6">
      {Object.entries(groupedDrinks).map(([date, drinks]) => (
        <div key={date} className="space-y-2">
          <h2 className="font-semibold">
            {format(new Date(date), "EEEE, MMMM d")}
          </h2>
          <div className="space-y-2">
            {drinks.map((drink) => (
              <DrinkLogItem key={drink.id} drink={drink} />
            ))}
          </div>
        </div>
      ))}

      {drinkLogs.length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          No drinks logged yet
        </div>
      )}
    </div>
  );
}

function DrinkLogItem({ drink }: { drink: DrinkLog }) {
  return (
    <div className="rounded-lg border p-4 hover:bg-muted/50">
      <div className="flex justify-between items-start">
        <div>
          <h3 className="font-medium">{drink.drink_name}</h3>
          <p className="text-sm text-muted-foreground">
            {drink.size_value}
            {drink.size_unit}, {drink.abv * 100}% ABV
          </p>
        </div>
        <div className="text-sm text-muted-foreground">
          {formatDistanceToNow(new Date(drink.logged_at), { addSuffix: true })}
        </div>
      </div>
    </div>
  );
}
