import type { Metadata } from "next";
import { IBM_Plex_Mono, Bricolage_Grotesque } from "next/font/google";
import { Toaster } from "@/components/ui/sonner";
import "./globals.css";
import { DrinkLogsProvider } from "@/contexts/drink-logs-context";

const fontSans = Bricolage_Grotesque({
  variable: "--font-bricolage-grotesque",
  subsets: ["latin"],
});

const fontMono = IBM_Plex_Mono({
  variable: "--font-ibm-plex-mono",
  subsets: ["latin"],
  weight: ["400", "700"],
});

export const metadata: Metadata = {
  title: "Sōber",
  description: "Sōber is a tool to help you take control of drinking habits.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${fontSans.variable} ${fontMono.variable} antialiased`}>
        <DrinkLogsProvider>
          {children}
          <Toaster />
        </DrinkLogsProvider>
      </body>
    </html>
  );
}
