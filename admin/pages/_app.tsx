import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Inter } from "next/font/google"; // Correct import for Next.js 13+
import Head from "next/head";
import { AuthProvider } from "@/hooks/useAuth";
import { Toaster } from "react-hot-toast"; // Optional: Good for notifications, implies I should install it or make a custom one. I'll skip installing for now and just use built-in or later add it.
// Actually, user requested "Error logging panel", etc. Toasts are good for "World Class UX".
// I'll stick to basic structure first.

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>TroyGH Admin - Ultimate Control System</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="TroyGH Admin Panel" />
        <meta name="robots" content="noindex, nofollow" />
        {/* Favicon requested from specific path, I can't access local user files outside workspace easily for build unless I copy them manually. 
            User said: "Add favicon from (C:\Users\SARTHAK\Downloads\G-Log-modified)".
            I cannot access that path. I will need to ask user to provide it or use a placeholder.
            I will set a placeholder link for now.
        */}
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={`${inter.variable} font-sans`}>
        <AuthProvider>
          <Component {...pageProps} />
          <Toaster
            position="top-right"
            toastOptions={{
              className: 'glass text-sm font-medium',
              style: {
                background: 'rgba(255, 255, 255, 0.8)',
                backdropFilter: 'blur(10px)',
                color: '#333',
              },
              success: {
                iconTheme: {
                  primary: '#10B981',
                  secondary: 'white',
                },
              }
            }}
          />
        </AuthProvider>
      </main>
    </>
  );
}
