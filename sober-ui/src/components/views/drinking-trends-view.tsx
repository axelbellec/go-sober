"use client";

import * as React from "react";
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { apiService } from "@/lib/api";

interface DrinkingTrendData {
  date: string;
  drink_count: number;
  total_standard_drinks: number;
}

const chartConfig = {
  drinks: {
    label: "Drinks",
  },
  drink_count: {
    label: "Number of Drinks",
    color: "hsl(var(--chart-1))",
  },
  total_standard_drinks: {
    label: "Standard Drinks",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;

export function DrinkingTrendsView() {
  const [timeRange, setTimeRange] = React.useState("90d");
  const [chartData, setChartData] = React.useState<DrinkingTrendData[]>([]);

  React.useEffect(() => {
    const fetchDrinkStats = async () => {
      try {
        const endDate = new Date().toISOString().split("T")[0];
        const startDate = new Date();
        startDate.setDate(startDate.getDate() - parseInt("90d"));

        const response = await apiService.getDrinkStats("daily", {
          start_date: startDate.toISOString().split("T")[0],
          end_date: endDate,
        });

        const formattedData = (() => {
          // Create a map of existing data points
          const dataMap = new Map(
            response.stats.map((item) => [
              item.time_period,
              {
                drink_count: item.drink_count,
                total_standard_drinks: item.total_standard_drinks,
              },
            ])
          );

          // Get the date range
          const dates = response.stats.map((item) => item.time_period);
          const startDate = new Date(
            Math.min(...dates.map((d) => new Date(d).getTime()))
          );
          const endDate = new Date(
            Math.max(...dates.map((d) => new Date(d).getTime()))
          );

          // Generate all dates in the range
          const allDates = [];
          const currentDate = new Date(startDate);
          while (currentDate <= endDate) {
            const dateStr = currentDate.toISOString().split("T")[0];
            allDates.push(dateStr);
            currentDate.setDate(currentDate.getDate() + 1);
          }

          // Create the final array with filled gaps
          return allDates.map((date) => ({
            date,
            drink_count: dataMap.get(date)?.drink_count ?? 0,
            total_standard_drinks:
              dataMap.get(date)?.total_standard_drinks ?? 0,
          }));
        })();

        setChartData(formattedData);
      } catch (error) {
        console.error("Failed to fetch drink stats:", error);
      }
    };

    fetchDrinkStats();
  }, []);

  const filteredData = chartData.filter((item) => {
    const date = new Date(item.date);
    const referenceDate = new Date().getDate();
    let daysToSubtract = 90;
    if (timeRange === "30d") {
      daysToSubtract = 30;
    } else if (timeRange === "7d") {
      daysToSubtract = 7;
    }
    const startDate = new Date();
    startDate.setDate(referenceDate - daysToSubtract);
    return date >= startDate;
  });

  return (
    <Card>
      <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
        <div className="grid flex-1 gap-1 text-center sm:text-left">
          <CardTitle>Drink Stats</CardTitle>
          <CardDescription>
            {timeRange === "90d"
              ? "Last 3 months"
              : timeRange === "30d"
              ? "Last 30 days"
              : "Last 7 days"}{" "}
            of drinking activity
          </CardDescription>
        </div>
        <Select value={timeRange} onValueChange={setTimeRange}>
          <SelectTrigger
            className="w-[160px] rounded-lg sm:ml-auto"
            aria-label="Select a value"
          >
            <SelectValue placeholder="Last 3 months" />
          </SelectTrigger>
          <SelectContent className="rounded-xl">
            <SelectItem value="90d" className="rounded-lg">
              Last 3 months
            </SelectItem>
            <SelectItem value="30d" className="rounded-lg">
              Last 30 days
            </SelectItem>
            <SelectItem value="7d" className="rounded-lg">
              Last 7 days
            </SelectItem>
          </SelectContent>
        </Select>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <AreaChart data={filteredData}>
            <defs>
              <linearGradient id="fillDrinkCount" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-drink-count)"
                  stopOpacity={0.75}
                />
                <stop
                  offset="98%"
                  stopColor="var(--color-drink-count)"
                  stopOpacity={0.01}
                />
              </linearGradient>
              <linearGradient
                id="fillStandardDrinks"
                x1="0"
                y1="0"
                x2="0"
                y2="1"
              >
                <stop
                  offset="5%"
                  stopColor="var(--color-standard-drinks)"
                  stopOpacity={0.75}
                />
                <stop
                  offset="98%"
                  stopColor="var(--color-standard-drinks)"
                  stopOpacity={0.01}
                />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => {
                const date = new Date(value);
                return date.toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                });
              }}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent indicator="dot" />}
            />
            <Area
              dataKey="total_standard_drinks"
              type="natural"
              fill="url(#fillStandardDrinks)"
              stroke="var(--color-standard-drinks)"
              stackId="a"
              baseValue={0}
            />
            <Area
              dataKey="drink_count"
              type="natural"
              fill="url(#fillDrinkCount)"
              stroke="var(--color-drink-count)"
              stackId="a"
              baseValue={0}
            />
            <ChartLegend content={<ChartLegendContent />} />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
