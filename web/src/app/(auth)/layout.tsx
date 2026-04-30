import Link from "next/link";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="relative flex min-h-screen items-center justify-center px-6 py-12">
      <div className="gradient-mesh absolute inset-0 -z-10" />
      <div className="w-full max-w-sm">
        <Link href="/" className="mb-8 flex items-center justify-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-md bg-[var(--accent)] text-[var(--accent-foreground)] text-sm font-bold">
            S
          </span>
          <span className="text-base font-semibold tracking-tight">Slugify</span>
        </Link>
        <div className="rounded-2xl border border-[var(--border)] bg-[var(--background)] p-6 shadow-xl shadow-indigo-500/5">
          {children}
        </div>
      </div>
    </div>
  );
}