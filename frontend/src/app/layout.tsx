import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "แชทอีสาน - หนองบัวลำภู | Isan Chat",
  description: "แอปพลิเคชันแชทเรียลไทม์ ธีมอีสานหนองบัวลำภู - Real-time Chat Application with Isan Nong Bua Lam Phu Theme",
  keywords: ["chat", "real-time", "isan", "nong bua lam phu", "หนองบัวลำภู", "แชท"],
  authors: [{ name: "Isan Chat Team" }],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="th">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Sarabun:wght@300;400;500;600;700&family=Prompt:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="antialiased">
        <div className="min-h-screen flex flex-col">
          {/* Header */}
          <header className="bg-gradient-to-r from-[var(--color-gold-500)] to-[var(--color-gold-600)] text-white shadow-lg">
            <div className="container mx-auto px-4 py-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="text-3xl">🏯</span>
                  <div>
                    <h1 className="text-xl font-bold">แชทอีสาน</h1>
                    <p className="text-xs text-[var(--color-gold-100)]">หนองบัวลำภู</p>
                  </div>
                </div>
                <nav className="flex items-center gap-4">
                  <span className="text-sm">ยินดีต้อนรับ! 🎋</span>
                </nav>
              </div>
            </div>
          </header>

          {/* Main Content */}
          <main className="flex-1 isan-pattern">
            {children}
          </main>

          {/* Footer */}
          <footer className="bg-[var(--color-earth-800)] text-[var(--color-earth-200)] py-4">
            <div className="container mx-auto px-4 text-center text-sm">
              <p>🎋 แชทอีสาน - หนองบัวลำภู | สร้างด้วย ❤️ จากดินแดนอีสาน</p>
              <p className="text-xs mt-1 text-[var(--color-earth-400)]">
                Next.js 16 • TailwindCSS 4 • Go Fiber • Redis 7.4 • PostgreSQL
              </p>
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}
