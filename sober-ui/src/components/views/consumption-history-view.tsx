"use client";

import { useEffect, useState } from "react";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { format, formatDistanceToNow } from "date-fns";
import { DrinkLog } from "@/lib/types/api";
import { DrinkLogForm } from "@/components/forms/drink-log-form";
import { Button } from "@/components/ui/button";
import { apiService } from "@/lib/api";
import { toast } from "sonner";

export function ConsumptionHistoryView() {
  const { drinkLogs, refreshDrinkLogs, fetchMoreDrinkLogs } = useDrinkLogs();
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    refreshDrinkLogs();
  }, [refreshDrinkLogs]);

  // Group drinks by date
  const groupedDrinks = (drinkLogs ?? []).reduce((groups, drink) => {
    // Safely parse the date and handle invalid dates
    const drinkDate = new Date(drink.logged_at);
    if (isNaN(drinkDate.getTime())) {
      return groups; // Skip invalid dates
    }

    const date = format(drinkDate, "yyyy-MM-dd");
    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(drink);
    groups[date].sort(
      (a, b) =>
        new Date(b.logged_at).getTime() - new Date(a.logged_at).getTime()
    );
    return groups;
  }, {} as Record<string, DrinkLog[]>);

  const handleLoadMore = async () => {
    setIsLoading(true);
    try {
      const newLogs = await fetchMoreDrinkLogs();
      setHasMore(newLogs.length > 0);
    } catch (error) {
      console.error("Error loading more drinks:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6 min-h-0">
      {Object.entries(groupedDrinks)
        .sort((a, b) => new Date(b[0]).getTime() - new Date(a[0]).getTime())
        .map(([date, drinks]) => (
          <div key={date} className="space-y-2">
            <h2 className="font-semibold sticky top-0 bg-background/95 backdrop-blur-sm py-2 z-10">
              {format(new Date(date), "EEEE, MMMM d")}
            </h2>
            <div className="space-y-2">
              {drinks.map((drink) => (
                <DrinkLogItem key={drink.id} drink={drink} />
              ))}
            </div>
          </div>
        ))}

      {hasMore && (drinkLogs ?? []).length > 0 && (
        <div className="flex justify-center py-4">
          <Button
            variant="outline"
            onClick={handleLoadMore}
            disabled={isLoading}
          >
            {isLoading ? "Loading..." : "Load More"}
          </Button>
        </div>
      )}

      {(drinkLogs ?? []).length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          No drinks logged yet
        </div>
      )}
    </div>
  );
}

function DrinkLogItem({ drink }: { drink: DrinkLog }) {
  const [isEditing, setIsEditing] = useState(false);
  const { refreshDrinkLogs } = useDrinkLogs();

  // Safely parse the date and handle invalid dates
  const loggedDate = new Date(drink.logged_at);
  const timeAgo = !isNaN(loggedDate.getTime())
    ? formatDistanceToNow(loggedDate, { addSuffix: true })
    : "Invalid date";

  if (isEditing) {
    return (
      <div className="rounded-lg border p-4">
        <DrinkLogForm
          initialDrinkLog={drink}
          onCancel={() => setIsEditing(false)}
          onDelete={async () => {
            try {
              await apiService.deleteDrinkLog(drink.id);
              await refreshDrinkLogs();
              toast.success("Drink log deleted");
            } catch (error) {
              console.error("Error deleting drink log:", error);
            }
          }}
          mode="edit"
        />
      </div>
    );
  }

  return (
    <div
      className="rounded-lg border p-4 hover:bg-muted/50"
      onClick={() => setIsEditing(true)}
      role="button"
      tabIndex={0}
    >
      <div className="flex justify-between items-start">
        <div>
          <h3 className="font-medium">{drink.name}</h3>
          <p className="text-sm text-muted-foreground">
            {drink.size_value}
            {drink.size_unit}, {(drink.abv * 100).toFixed(2)}% ABV
          </p>
        </div>
        <div className="text-sm text-muted-foreground">{timeAgo}</div>
      </div>
    </div>
  );
}
