"use client"

import type { FunctionComponent, PropsWithChildren } from "react"
import { ThemeProvider } from "./theme-provider.component"

export const Providers: FunctionComponent<PropsWithChildren> = ({ children }) => {
  return <ThemeProvider>{children}</ThemeProvider>
}
