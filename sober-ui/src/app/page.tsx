import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      {/* Hero Section */}
      <section className="flex flex-col items-center justify-center space-y-4 px-4 py-24 text-center md:py-32">
        <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl">
          Sōber 🧃
        </h1>
        <p className="max-w-[700px] text-lg text-muted-foreground sm:text-xl">
          Track your alcohol consumption intelligently. Make informed decisions
          about your drinking habits with real-time BAC monitoring.
        </p>
        <div className="flex flex-col gap-4 sm:flex-row">
          <Link href="/signup">
            <Button size="lg">Get Started</Button>
          </Link>
          <Link href="/login">
            <Button variant="outline" size="lg">
              Sign In
            </Button>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="container mx-auto px-4 py-16">
        <h2 className="mb-12 text-center text-3xl font-bold">Key Features</h2>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <Card>
            <CardHeader>
              <CardTitle>Real-time BAC Tracking</CardTitle>
              <CardDescription>
                Monitor your blood alcohol content in real-time with scientific
                accuracy
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                📊
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Drink History</CardTitle>
              <CardDescription>
                Keep a detailed log of your drinks and consumption patterns
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                📝
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Smart Insights</CardTitle>
              <CardDescription>
                Get personalized insights and recommendations based on your
                drinking patterns
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

      {/* CTA Section */}
      <section className="bg-muted py-16">
        <div className="container mx-auto px-4 text-center">
          <h2 className="mb-4 text-3xl font-bold">Ready to Take Control?</h2>
          <p className="mb-8 text-muted-foreground">
            Join thousands of users making smarter drinking decisions
          </p>
          <Link href="/signup">
            <Button size="lg" className="min-w-[200px]">
              Start Tracking Now
            </Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-6">
        <div className="container mx-auto px-4 text-center text-sm text-muted-foreground">
          <p>© 2024 Sōber. All rights reserved.</p>
          <p className="mt-2">
            Drink responsibly. This app is for informational purposes only.
          </p>
        </div>
      </footer>
    </div>
  );
}
