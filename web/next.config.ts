import type { NextConfig } from "next";

const backend =
  process.env.BACKEND_URL?.replace(/\/$/, "") ?? "http://localhost:9000";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/backend/:path*",
        destination: `${backend}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;