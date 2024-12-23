"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useState } from "react";

export default function Home() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <div className="flex min-h-screen flex-col">
      {/* Navigation Header */}
      <header className="border-b relative z-50">
        <div className="container mx-auto px-4">
          <nav className="flex h-16 items-center justify-between">
            <Link href="/" className="flex items-center space-x-2">
              <span className="text-xl font-bold">Sōber</span>
              <span className="text-2xl">🧃</span>
            </Link>

            {/* Desktop Navigation */}
            <div className="hidden space-x-8 md:flex">
              <Link
                href="/features"
                className="text-sm text-muted-foreground hover:text-primary"
              >
                Features
              </Link>
              <Link
                href="/pricing"
                className="text-sm text-muted-foreground hover:text-primary"
              >
                Pricing
              </Link>
              <Link
                href="/about"
                className="text-sm text-muted-foreground hover:text-primary"
              >
                About
              </Link>
              <Link
                href="/blog"
                className="text-sm text-muted-foreground hover:text-primary"
              >
                Blog
              </Link>
            </div>

            {/* Auth Buttons */}
            <div className="hidden space-x-4 md:flex">
              <Link href="/login">
                <Button variant="ghost" size="sm">
                  Sign In
                </Button>
              </Link>
              <Link href="/signup">
                <Button size="sm">Get Started</Button>
              </Link>
            </div>

            {/* Mobile Menu Button */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            >
              {isMobileMenuOpen ? (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              ) : (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <line x1="3" y1="12" x2="21" y2="12" />
                  <line x1="3" y1="6" x2="21" y2="6" />
                  <line x1="3" y1="18" x2="21" y2="18" />
                </svg>
              )}
            </Button>
          </nav>

          {/* Mobile Menu */}
          <div
            className={`${
              isMobileMenuOpen ? "flex" : "hidden"
            } md:hidden absolute left-0 right-0 top-16 z-50 flex-col bg-background border-b`}
          >
            <div className="flex flex-col space-y-4 p-4">
              <Link
                href="/features"
                className="text-sm text-muted-foreground hover:text-primary"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Features
              </Link>
              <Link
                href="/pricing"
                className="text-sm text-muted-foreground hover:text-primary"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Pricing
              </Link>
              <Link
                href="/about"
                className="text-sm text-muted-foreground hover:text-primary"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                About
              </Link>
              <Link
                href="/blog"
                className="text-sm text-muted-foreground hover:text-primary"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Blog
              </Link>
            </div>
            <div className="border-t p-4 space-y-2">
              <Link href="/login" className="w-full">
                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Sign In
                </Button>
              </Link>
              <Link href="/signup" className="w-full">
                <Button
                  size="sm"
                  className="w-full"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Get Started
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* Wrap all content below header in a div */}
      <div className={`${isMobileMenuOpen ? "blur-sm" : ""} transition-all`}>
        {/* Hero Section */}
        <section className="flex flex-col items-center justify-center space-y-4 px-4 py-24 text-center md:py-32">
          <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl">
            Mindful Drinking Made Simple 🧃
          </h1>
          <p className="max-w-[700px] text-lg text-muted-foreground sm:text-xl">
            Take control of your drinking habits with Sōber. Track your drinks,
            understand your limits, and make informed decisions about your
            alcohol consumption.
          </p>
          <div className="flex flex-col gap-4 sm:flex-row">
            <Link href="/signup">
              <Button size="lg">Start Your Journey</Button>
            </Link>
            <Link href="/login">
              <Button variant="outline" size="lg">
                Welcome Back
              </Button>
            </Link>
          </div>
        </section>

        {/* Features Section */}
        <section className="container mx-auto px-4 py-16">
          <h2 className="mb-12 text-center text-3xl font-bold">
            How Sōber Helps You
          </h2>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle>Smart BAC Tracking</CardTitle>
                <CardDescription>
                  Know where you stand with real-time blood alcohol tracking
                  that's both accurate and easy to understand
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  📊
                </div>
              </CardContent>
            </Card>

            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle>Personal Insights</CardTitle>
                <CardDescription>
                  Build better habits with a clear view of your drinking
                  patterns and personalized recommendations
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  📝
                </div>
              </CardContent>
            </Card>

            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle>Mindful Decisions</CardTitle>
                <CardDescription>
                  Make informed choices about your drinking with real-time
                  insights and gentle reminders
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  💡
                </div>
              </CardContent>
            </Card>
          </div>
        </section>

        {/* Benefits Section */}
        <section className="bg-muted/50 py-16">
          <div className="container mx-auto px-4">
            <h2 className="mb-12 text-center text-3xl font-bold">
              Your Journey with Sōber
            </h2>
            <div className="grid gap-8 md:grid-cols-2">
              <div className="flex items-start space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  🎯
                </div>
                <div>
                  <h3 className="text-xl font-semibold">
                    Science-Backed Tracking
                  </h3>
                  <p className="mt-2 text-muted-foreground">
                    Trust our precise BAC calculations to give you reliable,
                    actionable insights about your drinking.
                  </p>
                </div>
              </div>
              <div className="flex items-start space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  🔒
                </div>
                <div>
                  <h3 className="text-xl font-semibold">
                    Your Data, Your Privacy
                  </h3>
                  <p className="mt-2 text-muted-foreground">
                    Rest easy knowing your information is encrypted and
                    completely private. We never share your data.
                  </p>
                </div>
              </div>
              <div className="flex items-start space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  📱
                </div>
                <div>
                  <h3 className="text-xl font-semibold">Effortless Tracking</h3>
                  <p className="mt-2 text-muted-foreground">
                    Track your drinks in seconds with our intuitive interface
                    designed for real-world use.
                  </p>
                </div>
              </div>
              <div className="flex items-start space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  📈
                </div>
                <div>
                  <h3 className="text-xl font-semibold">Clear Insights</h3>
                  <p className="mt-2 text-muted-foreground">
                    Understand your habits with easy-to-read charts and
                    personalized insights that matter.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="bg-muted py-16">
          <div className="container mx-auto px-4 text-center">
            <h2 className="mb-4 text-3xl font-bold">
              Start Your Mindful Journey Today
            </h2>
            <p className="mb-8 text-muted-foreground">
              Join thousands making smarter choices about their drinking habits
            </p>
            <Link href="/signup">
              <Button size="lg" className="min-w-[200px]">
                Get Started Free
              </Button>
            </Link>
          </div>
        </section>

        {/* Footer */}
        <footer className="border-t py-12">
          <div className="container mx-auto px-4">
            <div className="grid gap-8 sm:grid-cols-2 md:grid-cols-4">
              {/* Company Info */}
              <div className="space-y-3">
                <h3 className="text-lg font-semibold">Sōber</h3>
                <p className="text-sm text-muted-foreground">
                  Making mindful drinking easier through technology.
                </p>
              </div>

              {/* Quick Links */}
              <div className="space-y-3">
                <h3 className="text-lg font-semibold">Quick Links</h3>
                <ul className="space-y-2">
                  <li>
                    <Link
                      href="/about"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      About Us
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/features"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Features
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/pricing"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Pricing
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/blog"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Blog
                    </Link>
                  </li>
                </ul>
              </div>

              {/* Support */}
              <div className="space-y-3">
                <h3 className="text-lg font-semibold">Support</h3>
                <ul className="space-y-2">
                  <li>
                    <Link
                      href="/help"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Help Center
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/contact"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Contact Us
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/faq"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      FAQ
                    </Link>
                  </li>
                </ul>
              </div>

              {/* Legal */}
              <div className="space-y-3">
                <h3 className="text-lg font-semibold">Legal</h3>
                <ul className="space-y-2">
                  <li>
                    <Link
                      href="/privacy"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Privacy Policy
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/terms"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Terms of Service
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/disclaimer"
                      className="text-sm text-muted-foreground hover:text-primary"
                    >
                      Disclaimer
                    </Link>
                  </li>
                </ul>
              </div>
            </div>

            <div className="mt-8 border-t pt-8 text-center text-sm text-muted-foreground">
              <p>© 2024 Sōber. All rights reserved.</p>
              <p className="mt-2">
                Drink responsibly. This app is for informational purposes only.
              </p>
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
}
