"use client"

import * as React from "react"

import { cn } from "@/lib/utils"

function Error({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="error"
      className={cn(
        "bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-neutral-100 px-4 py-3 rounded-md text-sm",
        className
      )}
      {...props}
    />
  )
}

export { Error }
