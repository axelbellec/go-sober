"use client";

import { useEffect, useState } from "react";
import { format } from "date-fns";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { BACCalculationResponse } from "@/lib/types/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { apiService } from "@/lib/api";

export function AlcoholLevelTimelineView() {
  const [data, setData] = useState<BACCalculationResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchBACTimeline() {
      try {
        // Now - 2hours
        const startTime = new Date(Date.now() - 2 * 60 * 60 * 1000);
        // Now + 2 hours
        const endTime = new Date(Date.now() + 4 * 60 * 60 * 1000);

        const data = await apiService.getBACTimeline(
          startTime.toISOString(),
          endTime.toISOString(),
          70, // Default weight, should be fetched from user profile
          "male", // Default gender, should be fetched from user profile
          2 // 2-minute intervals
        );
        setData(data);
      } catch (error) {
        console.error("Failed to fetch BAC timeline:", error);
      } finally {
        setIsLoading(false);
      }
    }

    fetchBACTimeline();
  }, []);

  if (isLoading) {
    return <div className="text-center">Loading...</div>;
  }

  if (!data) {
    return (
      <div className="text-center text-muted-foreground">No data available</div>
    );
  }

  const chartData = data.timeline.map((point) => ({
    time: format(new Date(point.time), "HH:mm"),
    bac: point.bac,
  }));

  const maxBAC = data.summary.max_bac;
  const drinkingSince = format(
    new Date(data.summary.drinking_since_time),
    "HH:mm"
  );
  const soberSince = data.summary.sober_since_time
    ? format(new Date(data.summary.sober_since_time), "HH:mm")
    : "Still drinking";

  return (
    <div className="w-full max-w-7xl mx-auto px-4 space-y-6">
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <Card className="shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Peak BAC</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">
              {(maxBAC * 100).toFixed(3)}%
            </div>
          </CardContent>
        </Card>
        <Card className="shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Started Drinking
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{drinkingSince}</div>
          </CardContent>
        </Card>
        <Card className="shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Sober Since</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{soberSince}</div>
          </CardContent>
        </Card>
      </div>

      <Card className="shadow-sm hover:shadow-md transition-shadow">
        <CardHeader>
          <CardTitle>BAC Timeline</CardTitle>
          <CardDescription>
            Your blood alcohol content over time
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] sm:h-[400px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart
                data={chartData}
                margin={{ top: 10, right: 10, left: 0, bottom: 0 }}
              >
                <defs>
                  <linearGradient id="bacGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop
                      offset="5%"
                      stopColor="hsl(var(--primary))"
                      stopOpacity={0.8}
                    />
                    <stop
                      offset="95%"
                      stopColor="hsl(var(--primary))"
                      stopOpacity={0.1}
                    />
                  </linearGradient>
                </defs>
                <CartesianGrid
                  vertical={false}
                  stroke="#888888"
                  opacity={0.2}
                />
                <XAxis
                  dataKey="time"
                  stroke="#888888"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  padding={{ left: 10, right: 10 }}
                />
                <YAxis
                  stroke="#888888"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={(value) => `${(value * 100).toFixed(2)}%`}
                />
                <Tooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      return (
                        <div className="rounded-lg border bg-background p-2 shadow-sm">
                          <div className="grid grid-cols-2 gap-2">
                            <div className="flex flex-col">
                              <span className="text-[0.70rem] uppercase text-muted-foreground">
                                Time
                              </span>
                              <span className="font-bold">
                                {payload[0].payload.time}
                              </span>
                            </div>
                            <div className="flex flex-col">
                              <span className="text-[0.70rem] uppercase text-muted-foreground">
                                BAC
                              </span>
                              <span className="font-bold">
                                {(Number(payload[0].value) * 100).toFixed(3)}%
                              </span>
                            </div>
                          </div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="bac"
                  stroke="hsl(var(--primary))"
                  fill="url(#bacGradient)"
                  fillOpacity={0.4}
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 6, fill: "hsl(var(--primary))" }}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
