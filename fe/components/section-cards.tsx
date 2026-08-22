import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Dashboard } from "@/lib/schemas"

export function SectionCards({ total_products }: Partial<Dashboard>) {
  return (
    <div className="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-2">
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Total Products</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {total_products}
          </CardTitle>
        </CardHeader>
      </Card>
    </div>
  )
}
