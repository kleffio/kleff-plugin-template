import { Card, CardHeader, CardTitle, CardContent } from "@kleff/ui";

export function ExamplePage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">My Plugin</h1>
        <p className="text-muted-foreground mt-1">
          This is a full page contributed by your plugin.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Getting started</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Edit <code>src/pages/ExamplePage.tsx</code> to build your page.
            You have access to all <code>@kleff/ui</code> components and the
            full Kleff design system.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
