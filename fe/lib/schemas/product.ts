import { z } from "zod"

export const ProductSchema = z.object({
      id: z.number(),
      sku: z.string().nullable(),
      name: z.string(),
      slug: z.string().nullable(),
      description: z.string().nullable(),
      price: z.number(),
      discount_price: z.number().nullable(),
      stock: z.number(),
      weight: z.number().nullable(),
      image_url: z.string().nullable(),
      status: z.string(),
      is_featured: z.boolean(),
})

export type Product = z.infer<typeof ProductSchema>
