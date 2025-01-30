"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import {
  DrinkTemplate,
  DrinkLog,
  ParseDrinkLogResponse,
} from "@/lib/types/api";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { apiService } from "@/lib/api";
import { useDrinkLogs } from "@/contexts/drink-logs-context";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";

const createDrinkLogSchema = z
  .object({
    drinkTemplateId: z.string().optional(),
    freeText: z.string().optional(),
    name: z.string().optional(),
    abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
    sizeValue: z.number().min(1, "Size must be greater than 0"),
    sizeUnit: z.enum(["cl", "ml"], {
      errorMap: () => ({ message: "Please select a valid unit" }),
    }),
  })
  .refine((data) => data.drinkTemplateId || data.freeText, {
    message: "Please either select a drink or describe what you had",
  });

const editDrinkLogSchema = z.object({
  name: z.string().min(1, "Name is required"),
  abv: z.number().min(0.01, "ABV must be greater than 0").max(100),
  sizeValue: z.number().min(1, "Size must be greater than 0"),
  sizeUnit: z.enum(["cl", "ml"], {
    errorMap: () => ({ message: "Please select a valid unit" }),
  }),
});

type DrinkLogFormValues = z.infer<typeof createDrinkLogSchema> &
  z.infer<typeof editDrinkLogSchema>;

interface DrinkLogFormProps {
  initialDrinkLog?: DrinkLog;
  onCancel?: () => void;
  mode?: "create" | "edit";
  onDelete?: () => void;
}

const drinkLogPlaceholders = [
  "e.g., had a glass of red wine",
  "e.g., drank a pint of beer",
  "e.g., had a mojito cocktail",
  "e.g., enjoyed a gin and tonic",
  "e.g., had a shot of tequila",
  "e.g., drank a vodka soda",
  "e.g., had a glass of prosecco",
  "e.g., enjoyed an old fashioned",
  "e.g., had a margarita",
  "e.g., drank a whiskey neat",
];

const getRandomPlaceholder = () => {
  return drinkLogPlaceholders[
    Math.floor(Math.random() * drinkLogPlaceholders.length)
  ];
};

interface DrinkDiff {
  name?: { from: string; to: string };
  size?: { from: number; to: number; unit: string };
  abv?: { from: number; to: number };
  unit?: { from: string; to: string };
}

function computeDrinkDiff(
  before: DrinkLog,
  after: {
    name: string;
    size_value: number;
    size_unit: string;
    abv: number;
  }
): DrinkDiff {
  const diff: DrinkDiff = {};

  if (before.name !== after.name) {
    diff.name = { from: before.name, to: after.name };
  }

  if (before.size_value !== after.size_value) {
    diff.size = {
      from: before.size_value,
      to: after.size_value,
      unit: after.size_unit,
    };
  }

  if (before.size_unit !== after.size_unit) {
    diff.unit = { from: before.size_unit, to: after.size_unit };
  }

  if (before.abv !== after.abv) {
    diff.abv = { from: before.abv * 100, to: after.abv * 100 };
  }

  return diff;
}

function formatDrinkDiff(diff: DrinkDiff): string {
  const changes: string[] = [];

  if (diff.name) {
    changes.push(`name from "${diff.name.from}" to "${diff.name.to}"`);
  }

  if (diff.size) {
    changes.push(
      `size from ${diff.size.from} to ${diff.size.to}${diff.size.unit}`
    );
  }

  if (diff.unit) {
    changes.push(`unit from ${diff.unit.from} to ${diff.unit.to}`);
  }

  if (diff.abv) {
    changes.push(
      `ABV from ${diff.abv.from.toFixed(1)}% to ${diff.abv.to.toFixed(1)}%`
    );
  }

  if (changes.length === 0) return "No changes made";
  if (changes.length === 1) return `Changed ${changes[0]}`;

  const lastChange = changes.pop();
  return `Changed ${changes.join(", ")} and ${lastChange}`;
}

export function DrinkLogForm({
  initialDrinkLog,
  onCancel,
  mode = "create",
  onDelete,
}: DrinkLogFormProps) {
  const [drinkTemplates, setDrinkTemplates] = useState<DrinkTemplate[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isParsingDrink, setIsParsingDrink] = useState(false);
  const [placeholder, setPlaceholder] = useState("Describe your drink...");
  const [showParsedForm, setShowParsedForm] = useState(false);

  useEffect(() => {
    setPlaceholder(getRandomPlaceholder());
  }, []);

  const form = useForm<DrinkLogFormValues>({
    resolver: zodResolver(
      mode === "create" ? createDrinkLogSchema : editDrinkLogSchema
    ),
    defaultValues: {
      drinkTemplateId:
        mode === "create" ? initialDrinkLog?.id?.toString() : undefined,
      freeText: "",
      name: mode === "edit" ? initialDrinkLog?.name ?? "" : undefined,
      abv: initialDrinkLog ? initialDrinkLog.abv * 100 : 0,
      sizeValue: initialDrinkLog?.size_value ?? 0,
      sizeUnit: (initialDrinkLog?.size_unit as "cl" | "ml") ?? "cl",
    },
    mode: "onChange",
  });

  const { refreshDrinkLogs } = useDrinkLogs();

  useEffect(() => {
    apiService
      .getDrinkTemplates()
      .then((data) => {
        const templates = Array.isArray(data.drink_templates)
          ? data.drink_templates
          : [];
        setDrinkTemplates(templates);
      })
      .catch((error) => {
        console.error("Failed to fetch drink templates:", error);
        setDrinkTemplates([]);
      });
  }, []);

  const onDrinkSelect = (drinkId: string) => {
    const drink = drinkTemplates.find((d) => d.id === parseInt(drinkId));
    if (drink) {
      form.setValue("freeText", "");
      form.setValue("drinkTemplateId", drink.id.toString());
      form.setValue("abv", drink.abv * 100, { shouldValidate: true });
      form.setValue("sizeValue", drink.size_value, { shouldValidate: true });
      form.setValue("sizeUnit", drink.size_unit as "cl" | "ml", {
        shouldValidate: true,
      });
    }
  };

  const parseDrink = async (value: string) => {
    if (!value) return;

    setIsParsingDrink(true);
    try {
      const response: ParseDrinkLogResponse = await apiService.parseDrinkLog(
        value
      );

      if (response.drink_parsed.success) {
        // Clear any existing template selection when using free text
        form.setValue("drinkTemplateId", undefined);

        // Set all the parsed values
        form.setValue("name", response.drink_parsed.name, {
          shouldValidate: true,
        });
        form.setValue("abv", response.drink_parsed.abv * 100, {
          shouldValidate: true,
        });
        form.setValue("sizeValue", response.drink_parsed.size_value, {
          shouldValidate: true,
        });
        form.setValue(
          "sizeUnit",
          response.drink_parsed.size_unit as "cl" | "ml",
          {
            shouldValidate: true,
          }
        );
        setShowParsedForm(true);
        toast.success("Drink details parsed successfully", {
          description: "Please review and adjust if needed",
          duration: 3000,
        });
      } else {
        console.warn(
          "Failed to parse drink:",
          response.drink_parsed.error_message
        );
        toast.error("Could not understand drink description", {
          description: response.drink_parsed.error_message,
          duration: 3000,
        });
      }
    } catch (error) {
      console.error("Failed to parse drink:", error);
      toast.error("Could not parse drink description", {
        description: "Please try describing your drink differently",
        duration: 3000,
      });
    } finally {
      setIsParsingDrink(false);
    }
  };

  const debouncedParseDrink = useDebounce(parseDrink, 1000);

  const onFreeTextChange = (value: string) => {
    form.setValue("freeText", value);
    setShowParsedForm(false);
    debouncedParseDrink(value);
  };

  async function onSubmit(data: DrinkLogFormValues) {
    if (!form.formState.isValid) return;

    setIsLoading(true);
    try {
      if (mode === "edit" && initialDrinkLog) {
        const updateData = {
          id: initialDrinkLog.id,
          name: data.name || "",
          type: initialDrinkLog.type || "custom",
          size_value: data.sizeValue,
          size_unit: data.sizeUnit,
          abv: data.abv / 100,
        };

        await apiService.updateDrinkLog(updateData);
        const diff = computeDrinkDiff(initialDrinkLog, {
          name: updateData.name,
          size_value: updateData.size_value,
          size_unit: updateData.size_unit,
          abv: updateData.abv,
        });
        toast.success(`Updated ${updateData.name}`, {
          description: formatDrinkDiff(diff),
          duration: 3000,
        });
      } else {
        const createData = (() => {
          // If using a template
          if (data.drinkTemplateId) {
            const selectedTemplate = drinkTemplates.find(
              (d) => d.id.toString() === data.drinkTemplateId
            );
            if (!selectedTemplate) {
              throw new Error("No drink template selected");
            }
            return {
              name: selectedTemplate.name,
              type: selectedTemplate.type,
              size_value: data.sizeValue,
              size_unit: data.sizeUnit,
              abv: data.abv / 100,
            };
          }
          // If using free text
          return {
            name: data.name || data.freeText || "Custom Drink",
            type: "custom",
            size_value: data.sizeValue,
            size_unit: data.sizeUnit,
            abv: data.abv / 100,
          };
        })();

        await apiService.createDrinkLog(createData);
        toast.success(`Added ${createData.name} to your log`, {
          description: `${createData.size_value}${createData.size_unit} at ${(
            createData.abv * 100
          ).toFixed(1)}% ABV`,
          duration: 3000,
        });
      }

      form.reset();
      refreshDrinkLogs();
      if (onCancel) onCancel();
      setShowParsedForm(false);
    } catch (error) {
      console.error(`Failed to ${mode} drink:`, error);
      toast.error(`Unable to ${mode === "edit" ? "update" : "log"} drink`, {
        description:
          "Something went wrong. Please check your input and try again.",
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-1">
        {mode === "create" ? (
          <Tabs defaultValue="template" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="template">Choose from templates</TabsTrigger>
              <TabsTrigger value="freetext">Describe your drink</TabsTrigger>
            </TabsList>

            <TabsContent value="template" className="space-y-4">
              <FormField
                control={form.control}
                name="drinkTemplateId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Select Drink</FormLabel>
                    <Select
                      onValueChange={(value) => {
                        form.setValue("freeText", "");
                        onDrinkSelect(value);
                      }}
                      value={field.value?.toString() || ""}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select a drink" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {drinkTemplates.map((drink) => (
                          <SelectItem
                            key={drink.id}
                            value={drink.id.toString()}
                          >
                            {drink.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="space-y-4">
                <FormField
                  control={form.control}
                  name="abv"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>ABV %</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step="0.1"
                          placeholder="Alcohol percentage"
                          {...field}
                          value={field.value || ""}
                          onChange={(e) =>
                            field.onChange(
                              e.target.value ? parseFloat(e.target.value) : ""
                            )
                          }
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex gap-4">
                  <FormField
                    control={form.control}
                    name="sizeValue"
                    render={({ field }) => (
                      <FormItem className="flex-1">
                        <FormLabel>Size</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            placeholder="Amount"
                            {...field}
                            value={field.value || ""}
                            onChange={(e) =>
                              field.onChange(
                                e.target.value ? parseInt(e.target.value) : ""
                              )
                            }
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="sizeUnit"
                    render={({ field }) => (
                      <FormItem className="w-24">
                        <FormLabel>Unit</FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Unit" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectItem value="cl">cl</SelectItem>
                            <SelectItem value="ml">ml</SelectItem>
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>
            </TabsContent>

            <TabsContent value="freetext" className="space-y-4">
              <FormField
                control={form.control}
                name="freeText"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Describe your drink</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          placeholder={placeholder}
                          autoComplete="off"
                          {...field}
                          onChange={(e) => onFreeTextChange(e.target.value)}
                        />
                        {isParsingDrink && (
                          <div className="absolute right-3 top-1/2 -translate-y-1/2">
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                          </div>
                        )}
                      </div>
                    </FormControl>
                    <FormDescription>
                      Describe what you had and we'll try to understand it
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {showParsedForm && (
                <div className="space-y-4 rounded-lg border p-4">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Drink Name</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="abv"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>ABV %</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            step="0.1"
                            placeholder="Alcohol percentage"
                            {...field}
                            value={field.value || ""}
                            onChange={(e) =>
                              field.onChange(
                                e.target.value ? parseFloat(e.target.value) : ""
                              )
                            }
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-4">
                    <FormField
                      control={form.control}
                      name="sizeValue"
                      render={({ field }) => (
                        <FormItem className="flex-1">
                          <FormLabel>Size</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder="Amount"
                              {...field}
                              value={field.value || ""}
                              onChange={(e) =>
                                field.onChange(
                                  e.target.value ? parseInt(e.target.value) : ""
                                )
                              }
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="sizeUnit"
                      render={({ field }) => (
                        <FormItem className="w-24">
                          <FormLabel>Unit</FormLabel>
                          <Select
                            onValueChange={field.onChange}
                            value={field.value}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Unit" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="cl">cl</SelectItem>
                              <SelectItem value="ml">ml</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </div>
              )}
            </TabsContent>
          </Tabs>
        ) : (
          <div className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Drink Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="abv"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ABV %</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.1"
                      placeholder="Alcohol percentage"
                      {...field}
                      value={field.value || ""}
                      onChange={(e) =>
                        field.onChange(
                          e.target.value ? parseFloat(e.target.value) : ""
                        )
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex gap-4">
              <FormField
                control={form.control}
                name="sizeValue"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Size</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        placeholder="Amount"
                        {...field}
                        value={field.value || ""}
                        onChange={(e) =>
                          field.onChange(
                            e.target.value ? parseInt(e.target.value) : ""
                          )
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="sizeUnit"
                render={({ field }) => (
                  <FormItem className="w-24">
                    <FormLabel>Unit</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Unit" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="cl">cl</SelectItem>
                        <SelectItem value="ml">ml</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </div>
        )}

        <div className="flex gap-2">
          <Button
            type="submit"
            className="flex-1"
            disabled={isLoading || !form.formState.isValid}
          >
            {isLoading
              ? `${mode === "edit" ? "Updating..." : "Logging..."}`
              : `${mode === "edit" ? "Update" : "Log"} Drink`}
          </Button>
          {mode === "edit" && (
            <>
              {onCancel && (
                <Button type="button" variant="outline" onClick={onCancel}>
                  Cancel
                </Button>
              )}
              {onDelete && (
                <Button
                  type="button"
                  variant="destructive"
                  onClick={async () => {
                    if (initialDrinkLog && onDelete) {
                      onDelete();
                    }
                  }}
                >
                  Delete
                </Button>
              )}
            </>
          )}
        </div>
      </form>
    </Form>
  );
}
