import type { NextConfig } from "next";

const backend =
  process.env.BACKEND_URL?.replace(/\/$/, "") ?? "http://localhost:9000";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      { source: "/backend/:path*", destination: `${backend}/api/:path*` },
      { source: "/auth-api/:path*", destination: `${backend}/auth/:path*` },
    ];
  },
};

export default nextConfig;