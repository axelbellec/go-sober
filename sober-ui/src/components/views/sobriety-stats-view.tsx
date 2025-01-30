"use client";

import {
  Label,
  PolarGrid,
  PolarRadiusAxis,
  RadialBar,
  RadialBarChart,
} from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ChartConfig, ChartContainer } from "@/components/ui/chart";
import { apiService } from "@/lib/api";
import { useState, useEffect } from "react";
import type { MonthlyBACStats } from "@/lib/types/api";

const chartConfig = {
  sober: {
    label: "Sober Days",
    color: "hsl(var(--success))",
  },
} satisfies ChartConfig;

const transformDataForChart = (stats: MonthlyBACStats[]) => {
  if (!stats.length)
    return [
      {
        name: "Sober Days",
        value: 0,
        fill: "hsl(var(--success))",
      },
    ];
  const latest = stats[stats.length - 1];
  const total = latest.total || 1;

  return [
    {
      name: "Sober Days",
      value: (latest.counts.sober / total) * 100,
      fill: "hsl(var(--success))",
    },
  ];
};

export function SobrietyStatsView() {
  const [data, setData] = useState<{ stats: MonthlyBACStats[] } | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const endDate = new Date().toISOString().split("T")[0];
        const startDate = new Date();
        startDate.setMonth(startDate.getMonth() - 1); // Only need current month

        const response = await apiService.getMonthlyBACStats({
          start_date: startDate.toISOString().split("T")[0],
          end_date: endDate,
        });

        setData(response);
      } catch (error) {
        console.error("Failed to fetch monthly BAC stats:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const chartData = data ? transformDataForChart(data.stats) : [];
  const soberDays = data?.stats[data.stats.length - 1]?.counts.sober ?? 0;
  const monthYear = data?.stats[data.stats.length - 1]
    ? `${new Date(
        data.stats[data.stats.length - 1].year,
        data.stats[data.stats.length - 1].month - 1
      ).toLocaleString("default", { month: "long", year: "numeric" })}`
    : "";

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Sober Days</CardTitle>
        <CardDescription>{monthYear}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <RadialBarChart
            data={chartData}
            startAngle={0}
            endAngle={360}
            innerRadius={80}
            outerRadius={110}
          >
            <PolarGrid
              gridType="circle"
              radialLines={false}
              stroke="none"
              className="first:fill-muted last:fill-background"
              polarRadius={[86, 74]}
            />
            <RadialBar dataKey="value" background cornerRadius={10} />
            <PolarRadiusAxis tick={false} tickLine={false} axisLine={false}>
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="fill-foreground text-4xl font-bold"
                        >
                          {soberDays}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground"
                        >
                          Days
                        </tspan>
                      </text>
                    );
                  }
                }}
              />
            </PolarRadiusAxis>
          </RadialBarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="flex items-center gap-2 font-medium leading-none">
          {!isLoading && `${chartData[0]?.value.toFixed(1)}% of the month`}
        </div>
        <div className="leading-none text-muted-foreground">
          Days completely sober this month
        </div>
      </CardFooter>
    </Card>
  );
}
