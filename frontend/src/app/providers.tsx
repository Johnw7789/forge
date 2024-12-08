"use client";
import * as React from "react";
import { NextUIProvider } from "@nextui-org/system";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import { ThemeProviderProps } from "next-themes/dist/types";
import { Layout } from "../components/layout/layout";
import {useRouter} from 'next/router';
import { RecoilRoot } from "recoil";

export interface ProvidersProps {
  children: React.ReactNode;
  themeProps?: ThemeProviderProps;
}

export function Providers({ children, themeProps }: ProvidersProps) {
  // const router = useRouter();

  return (
    <RecoilRoot>
    <NextUIProvider >
      <NextThemesProvider defaultTheme="dark" attribute="class" {...themeProps}>
        <Layout>
          {children}
        </Layout>
      </NextThemesProvider>
    </NextUIProvider>
    </RecoilRoot>
  );
}
