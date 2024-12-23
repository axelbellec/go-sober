interface PageLayoutProps {
  children: React.ReactNode;
  heading: string;
  subheading?: string;
  className?: string;
}

export function PageLayout({
  children,
  heading,
  subheading,
  className,
}: PageLayoutProps) {
  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-6">
        <div className={`mx-auto flex w-full flex-col space-y-6 ${className}`}>
          <div className="flex flex-col space-y-2 text-center">
            <h1 className="text-2xl md:text-2xl font-semibold tracking-tight">
              {heading}
            </h1>
            {subheading && (
              <p className="text-base md:text-sm text-muted-foreground">
                {subheading}
              </p>
            )}
          </div>
          {children}
        </div>
      </div>
    </div>
  );
}
