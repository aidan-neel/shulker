import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { Slot } from "radix-ui"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
    [
        "group/button relative inline-flex h-max shrink-0 cursor-pointer items-center justify-center whitespace-nowrap rounded-lg outline-brand transition duration-100 ease-linear before:absolute select-none",
        "focus-visible:outline-2 focus-visible:outline-offset-2",
        "in-data-input-wrapper:shadow-xs in-data-input-wrapper:focus:!z-50",
        "in-data-input-wrapper:in-data-leading:-mr-px in-data-input-wrapper:in-data-leading:rounded-r-none in-data-input-wrapper:in-data-leading:before:rounded-r-none",
        "in-data-input-wrapper:in-data-trailing:-ml-px in-data-input-wrapper:in-data-trailing:rounded-l-none in-data-input-wrapper:in-data-trailing:before:rounded-l-none",
        "disabled:cursor-not-allowed disabled:opacity-50 in-data-input-wrapper:disabled:opacity-100",
        "*:data-icon:pointer-events-none *:data-icon:size-5 *:data-icon:shrink-0 *:data-icon:transition-inherit-all",
        "[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
    ].join(" "),
    {
        variants: {
            variant: {
                primary: [
                    "bg-primary text-primary-foreground shadow-xs ring-1 ring-transparent ring-inset",
                    "hover:bg-primary/90 data-[loading=true]:bg-primary/90",
                    "before:absolute before:inset-px before:border before:border-white/12 before:mask-b-from-0% before:rounded-[7px]",
                    "*:data-icon:text-primary-foreground/60 hover:*:data-icon:text-primary-foreground/70",
                ].join(" "),
                secondary: [
                    "bg-background text-secondary-foreground shadow-xs ring-1 ring-border ring-inset",
                    "hover:bg-secondary/80 data-[loading=true]:bg-secondary/80",
                    "*:data-icon:text-muted-foreground hover:*:data-icon:text-muted-foreground/80",
                ].join(" "),
                tertiary: [
                    "text-muted-foreground hover:bg-accent hover:text-accent-foreground data-[loading=true]:bg-accent",
                    "*:data-icon:text-muted-foreground hover:*:data-icon:text-muted-foreground/80",
                ].join(" "),
                "link-color": [
                    "justify-normal rounded p-0! text-primary hover:text-primary/80",
                    "*:data-text:underline *:data-text:decoration-transparent hover:*:data-text:decoration-primary/50",
                    "*:data-icon:text-primary/70 hover:*:data-icon:text-primary/90",
                ].join(" "),
                "link-gray": [
                    "justify-normal rounded p-0! text-muted-foreground hover:text-foreground",
                    "*:data-text:underline *:data-text:decoration-transparent hover:*:data-text:decoration-muted-foreground",
                    "*:data-icon:text-muted-foreground hover:*:data-icon:text-foreground",
                ].join(" "),
                "primary-destructive": [
                    "bg-destructive text-white shadow-xs ring-1 ring-transparent ring-inset outline-destructive",
                    "hover:bg-destructive/90 data-[loading=true]:bg-destructive/90",
                    "before:absolute before:inset-px before:border before:border-white/12 before:mask-b-from-0% before:rounded-[7px]",
                    "*:data-icon:text-white/60 hover:*:data-icon:text-white/70",
                ].join(" "),
                "secondary-destructive": [
                    "bg-background text-destructive shadow-xs ring-1 ring-destructive/30 ring-inset outline-destructive",
                    "hover:bg-destructive/10 data-[loading=true]:bg-destructive/10",
                    "*:data-icon:text-destructive/70 hover:*:data-icon:text-destructive",
                ].join(" "),
                "tertiary-destructive": [
                    "text-destructive outline-destructive hover:bg-destructive/10 data-[loading=true]:bg-destructive/10",
                    "*:data-icon:text-destructive/70 hover:*:data-icon:text-destructive",
                ].join(" "),
                "link-destructive": [
                    "justify-normal rounded p-0! text-destructive outline-destructive hover:text-destructive/80",
                    "*:data-text:underline *:data-text:decoration-transparent *:data-text:underline-offset-2 hover:*:data-text:decoration-current",
                    "*:data-icon:text-destructive/70 hover:*:data-icon:text-destructive",
                ].join(" "),
            },
            size: {
                xs: [
                    "gap-1 rounded-md px-2.5 py-1.5 text-xs font-semibold before:rounded-[7px] data-[icon-only=true]:p-2",
                    "in-data-input-wrapper:px-3.5 in-data-input-wrapper:py-2.5 in-data-input-wrapper:data-[icon-only=true]:p-2.5",
                    "*:data-icon:size-3.5 *:data-icon:stroke-[2.25px]",
                ].join(" "),
                sm: [
                    "gap-1 rounded-md px-3 py-2 text-sm font-semibold before:rounded-[7px] data-[icon-only=true]:p-2",
                    "in-data-input-wrapper:px-3.5 in-data-input-wrapper:py-2.5 in-data-input-wrapper:data-[icon-only=true]:p-2.5",
                ].join(" "),
                md: [
                    "gap-1 rounded-md px-3.5 py-2.5 text-sm font-semibold before:rounded-[7px] data-[icon-only=true]:p-2.5",
                    "in-data-input-wrapper:gap-1.5 in-data-input-wrapper:px-4 in-data-input-wrapper:data-[icon-only=true]:p-3",
                ].join(" "),
                lg: "gap-1.5 rounded-xl px-4 py-2.5 text-base font-semibold before:rounded-[calc(0.75rem-1px)] data-[icon-only=true]:p-3",
                xl: "gap-1.5 rounded-xl px-4.5 py-3 text-base font-semibold before:rounded-[calc(0.75rem-1px)] data-[icon-only=true]:p-3.5",
            },
        },
        defaultVariants: {
            variant: "primary",
            size: "md",
        },
    }
)

function Button({
    className,
    variant = "primary",
    size = "md",
    asChild = false,
    isLoading,
    isDisabled,
    iconLeading: IconLeading,
    iconTrailing: IconTrailing,
    showTextWhileLoading,
    children,
    ...props
}: React.ComponentProps<"button"> &
    VariantProps<typeof buttonVariants> & {
        asChild?: boolean
        isLoading?: boolean
        isDisabled?: boolean
        iconLeading?: React.FC<{ className?: string }> | React.ReactNode
        iconTrailing?: React.FC<{ className?: string }> | React.ReactNode
        showTextWhileLoading?: boolean
    }) {
    const Comp = asChild ? Slot.Root : "button"

    const isLinkType = ["link-gray", "link-color", "link-destructive"].includes(variant ?? "")
    const isIconOnly = (IconLeading || IconTrailing) && !children

    const iconClass = "pointer-events-none size-5 shrink-0 transition-inherit-all"

    const renderIcon = (icon: React.FC<{ className?: string }> | React.ReactNode, dataAttr: string) => {
        if (!icon) return null
        if (React.isValidElement(icon)) return icon
        if (typeof icon === "function") {
            const IconComp = icon as React.FC<{ className?: string }>
            return <IconComp data-icon={dataAttr} className={iconClass} />
        }
        return null
    }

    return (
        <Comp
            data-slot="button"
            data-variant={variant}
            data-size={size}
            data-loading={isLoading ? true : undefined}
            data-icon-only={isIconOnly ? true : undefined}
            disabled={isDisabled || isLoading}
            className={cn(
                buttonVariants({ variant, size }),
                isLinkType && (size === "xs" || size === "sm"
                    ? "gap-1 *:data-text:underline-offset-3"
                    : "gap-1.5 *:data-text:underline-offset-4"),
                isLoading && (showTextWhileLoading
                    ? "[&>*:not([data-icon=loading]):not([data-text])]:hidden pointer-events-none"
                    : "[&>*:not([data-icon=loading])]:invisible pointer-events-none"),
                className,
            )}
            {...props}
        >
            {renderIcon(IconLeading, "leading")}

            {isLoading && (
                <svg
                    fill="none"
                    data-icon="loading"
                    viewBox="0 0 20 20"
                    className={cn(
                        iconClass,
                        !showTextWhileLoading && "absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"
                    )}
                >
                    <circle className="stroke-current opacity-30" cx="10" cy="10" r="8" fill="none" strokeWidth="2" />
                    <circle
                        className="origin-center animate-spin stroke-current"
                        cx="10" cy="10" r="8"
                        fill="none"
                        strokeWidth="2"
                        strokeDasharray="12.5 50"
                        strokeLinecap="round"
                    />
                </svg>
            )}

            {children && (
                <span
                    data-text
                    className={cn("transition-inherit-all", !isLinkType && "px-0.5")}
                >
                    {children}
                </span>
            )}

            {renderIcon(IconTrailing, "trailing")}
        </Comp>
    )
}

export { Button, buttonVariants }
