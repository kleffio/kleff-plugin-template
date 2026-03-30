import { Card, CardHeader, CardTitle, CardContent } from "@kleff/ui";

/**
 * ExampleWidget is injected into the dashboard via the "dashboard.widgets" slot.
 * Keep widgets compact — they share space with other dashboard cards.
 */
export function ExampleWidget() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm font-medium">My Plugin Widget</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-2xl font-bold">—</p>
        <p className="text-xs text-muted-foreground mt-1">
          Edit <code>src/components/ExampleWidget.tsx</code>
        </p>
      </CardContent>
    </Card>
  );
}
