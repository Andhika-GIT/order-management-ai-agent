'use client'

import { useState } from "react"
import { Product } from "@/lib/schemas"
import { Checkbox } from "@/components/ui/checkbox"
import { ColumnDef } from "@tanstack/react-table"
import {
    DropdownMenu, DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { IconDotsVertical } from "@tabler/icons-react"
import { UploadProductImageModal } from "@/components/molecules"
import { Badge } from "@/components/ui/badge"

const ProductActionsCell: React.FC<{ product: Product }> = ({ product }) => {
  const [open, setOpen] = useState(false)

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
            size="icon"
          >
            <IconDotsVertical />
            <span className="sr-only">Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-40">
          <DropdownMenuItem>Edit</DropdownMenuItem>
          <DropdownMenuItem>Make a copy</DropdownMenuItem>
          <DropdownMenuItem>Favorite</DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={() => setOpen(true)}>Upload Image</DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem variant="destructive">Delete</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <UploadProductImageModal productId={product.id} open={open} onOpenChange={setOpen} />
    </>
  )
}

export const columns: ColumnDef<Product>[] = [

  {
    id: "select",
    header: ({ table }) => (
      <div className="flex items-center justify-center">
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      </div>
    ),
    cell: ({ row }) => (
      <div className="flex items-center justify-center">
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      </div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => {
        return row.original.name
    },
  },
  {
    accessorKey: "sku",
    header: "SKU",
    cell: ({ row }) => {
        return row.original.sku ?? "-"
    },
  },
  {
    accessorKey: "price",
    header: "Price",
    cell: ({ row }) => {
        return row.original.price
    },
  },
  {
    accessorKey: "stock",
    header: "Stock",
    cell: ({ row }) => {
        return row.original.stock
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
        return <Badge variant="outline" className="text-muted-foreground px-1.5">{row.original.status}</Badge>
    },
  },
  {
    accessorKey: "is_featured",
    header: "Featured",
    cell: ({ row }) => {
        return row.original.is_featured ? "Yes" : "No"
    },
  },

  {
    id: "actions",
    cell: ({ row }) => <ProductActionsCell product={row.original} />,
  },
]
