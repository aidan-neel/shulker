"use client";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AnimatePresence, motion } from "motion/react"
import { Spinner } from "@/components/ui/spinner";
import { authClient } from "@/lib/rpc/clients";
import { useRouter } from "next/navigation"
import { toast } from "sonner";
import { z } from "zod";

const schema = z.object({
    email: z.string().email("Invalid email address"),
    password: z.string().min(8, "Password must be at least 8 characters"),
})

type FieldErrors = Partial<Record<keyof z.infer<typeof schema>, string>>

export default function Register() {
    const router = useRouter()

    const [isLoading, setIsLoading] = useState(false)
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})

    async function register() {
        setFieldErrors({})

        const result = schema.safeParse({ email, password })
        if (!result.success) {
            const errors: FieldErrors = {}
            for (const issue of result.error.issues) {
                const field = issue.path[0] as keyof FieldErrors
                errors[field] = issue.message
            }
            setFieldErrors(errors)
            return
        }

        setIsLoading(true)
        try {
            await authClient.register({ email, password })
            toast.success("Successfully registered!")
            router.push('/storage')
        } catch (err) {
            console.error(err)
            toast.error("Something went wrong. Please try again.")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50 font-sans dark:bg-black">
            <div className="p-12 py-24 bg-background gap-8 flex items-center justify-center w-[27rem] flex-col rounded-lg">
                <header className="flex flex-col items-center gap-1">
                    <h1 className="text-2xl font-semibold font-heading">
                        Create an account
                    </h1>
                    <h2 className="text-muted-foreground">
                        Start saving all the files you need.
                    </h2>
                </header>

                <fieldset
                    disabled={isLoading}
                    className="w-full flex flex-col gap-5 disabled:opacity-50 disabled:pointer-events-none transition-opacity duration-200"
                >
                    <Label htmlFor="email" className="flex flex-col w-full items-start gap-2">
                        Email
                        <Input
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            id="email"
                            placeholder="Enter your email"
                            type="text"
                            aria-invalid={!!fieldErrors.email}
                            className={fieldErrors.email ? "border-destructive focus-visible:ring-destructive" : ""}
                        />
                        {fieldErrors.email && (
                            <span className="text-destructive text-xs">{fieldErrors.email}</span>
                        )}
                    </Label>
                    <Label htmlFor="password" className="flex flex-col w-full items-start gap-2">
                        Password
                        <Input
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            id="password"
                            placeholder="Create a password"
                            type="password"
                            aria-invalid={!!fieldErrors.password}
                            className={fieldErrors.password ? "border-destructive focus-visible:ring-destructive" : ""}
                        />
                        {fieldErrors.password && (
                            <span className="text-destructive text-xs">{fieldErrors.password}</span>
                        )}
                    </Label>

                    <Button onClick={register} className="overflow-hidden">
                        <div className="flex flex-row items-center gap-1.5">
                            <AnimatePresence initial={false} mode="popLayout">
                                {isLoading && (
                                    <motion.div
                                        key="spinner"
                                        initial={{ opacity: 0, scale: 0, width: 0 }}
                                        animate={{ opacity: 1, scale: 1, width: 16 }}
                                        exit={{ opacity: 0, scale: 0, width: 0 }}
                                        transition={{ duration: 0.2, ease: [0.4, 0, 0.2, 1] }}
                                        className="shrink-0"
                                    >
                                        <Spinner className="size-4" />
                                    </motion.div>
                                )}
                            </AnimatePresence>
                            <span>Get Started</span>
                        </div>
                    </Button>
                </fieldset>
            </div>
        </div>
    );
}
