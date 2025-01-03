"use client";

import { TrendingUp } from "lucide-react";
import { PolarAngleAxis, PolarGrid, Radar, RadarChart } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { apiService } from "@/lib/api";
import { useState, useEffect } from "react";
import type { MonthlyBACStats } from "@/lib/types/api";

const chartConfig = {
  sober: {
    label: "Sober Days",
    color: "hsl(var(--chart-1))",
  },
  light: {
    label: "Light Drinking",
    color: "hsl(var(--chart-2))",
  },
  heavy: {
    label: "Heavy Drinking",
    color: "hsl(var(--chart-3))",
  },
} satisfies ChartConfig;

export function SobrietyDashboardView() {
  const [chartData, setChartData] = useState<any[]>([]);
  const [trend, setTrend] = useState<number | null>(null);

  useEffect(() => {
    const fetchMonthlyStats = async () => {
      try {
        const endDate = new Date().toISOString().split("T")[0];
        const startDate = new Date();
        startDate.setMonth(startDate.getMonth() - 6);

        const response = await apiService.getMonthlyBACStats({
          start_date: startDate.toISOString().split("T")[0],
          end_date: endDate,
        });

        // Transform API response to match chart data format
        const formattedData = response.stats.map((item: MonthlyBACStats) => {
          // Calculate total days for the month

          return {
            month: new Date(item.year, item.month - 1).toLocaleString(
              "default",
              {
                month: "long",
                year: "numeric",
              }
            ),
            sober: item.counts.sober || 0,
            light: item.counts.light || 0,
            heavy: item.counts.heavy || 0,
            soberPercentage: (item.counts.sober / item.total) * 100,
          };
        });

        setChartData(formattedData);

        // Calculate trend using percentages instead of raw counts
        if (formattedData.length >= 2) {
          const lastMonth = formattedData[formattedData.length - 1];
          const previousMonth = formattedData[formattedData.length - 2];
          const trendPercentage =
            lastMonth.soberPercentage - previousMonth.soberPercentage;
          setTrend(trendPercentage);
        }
      } catch (error) {
        console.error("Failed to fetch monthly BAC stats:", error);
      }
    };

    fetchMonthlyStats();
  }, []);

  return (
    <Card>
      <CardHeader className="items-center pb-4">
        <CardTitle>Sober Days</CardTitle>
        <CardDescription>
          Showing total sober days for the last 6 months
        </CardDescription>
      </CardHeader>
      <CardContent className="pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <RadarChart
            data={chartData}
            margin={{
              top: 10,
              right: 10,
              bottom: 10,
              left: 10,
            }}
          >
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent indicator="line" />}
            />
            <PolarAngleAxis
              dataKey="month"
              tick={({ x, y, textAnchor, value, index, ...props }) => {
                const data = chartData[index];

                return (
                  <text
                    x={x}
                    y={index === 0 ? y - 10 : y}
                    textAnchor={textAnchor}
                    fontSize={13}
                    fontWeight={500}
                    {...props}
                  >
                    <tspan>{(data?.soberPercentage || 0).toFixed(0)}%</tspan>
                    <tspan
                      x={x}
                      dy={"1rem"}
                      fontSize={12}
                      className="fill-muted-foreground"
                    >
                      {data?.month}
                    </tspan>
                  </text>
                );
              }}
            />

            <PolarGrid />
            <Radar
              dataKey="sober"
              fill="var(--color-sober)"
              fillOpacity={0.3}
              stroke="var(--color-sober)"
              strokeOpacity={0.8}
            />
            <Radar
              dataKey="light"
              fill="var(--color-light)"
              fillOpacity={0.3}
              stroke="var(--color-light)"
              strokeOpacity={0.8}
            />
            <Radar
              dataKey="heavy"
              fill="var(--color-heavy)"
              fillOpacity={0.3}
              stroke="var(--color-heavy)"
              strokeOpacity={0.8}
            />
          </RadarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="flex justify-center gap-4 mb-2">
          {Object.entries(chartConfig).map(([key, config]) => (
            <div key={key} className="flex items-center gap-1">
              <div
                className="w-3 h-3 rounded"
                style={{ backgroundColor: config.color }}
              />
              <span className="text-sm">{config.label}</span>
            </div>
          ))}
        </div>
        <div className="flex items-center gap-2 font-medium leading-none">
          {trend !== null && (
            <>
              {trend > 0
                ? `Great progress! You're ${Math.abs(trend).toFixed(
                    1
                  )}% more sober this month`
                : `${Math.abs(trend).toFixed(
                    1
                  )}% fewer sober days this month - keep going!`}
              <TrendingUp
                className={`h-4 w-4 ${trend < 0 ? "rotate-180" : ""}`}
              />
            </>
          )}
        </div>
        <div className="flex items-center gap-2 leading-none text-muted-foreground">
          {chartData.length > 0
            ? `${chartData[0].month} - ${chartData[chartData.length - 1].month}`
            : "No data"}
        </div>
      </CardFooter>
    </Card>
  );
}
