using Yarp.ReverseProxy.Configuration;

namespace gateway.YarpUtils;

public static class Routes
{
    public static IReadOnlyList<RouteConfig> GetRoutes()
    {
        return
        [
            new RouteConfig
            {
                RouteId = "auth-public",
                CorsPolicy = "WebOriginCorsPolicy",
                ClusterId = "auth-cluster",
                Match = new RouteMatch { Path = "/api/auth/{**catch-all}" },
                Transforms = [new Dictionary<string, string> { { "PathRemovePrefix", "/api/auth" } }],
                AuthorizationPolicy = "anonymous",
            },
            new RouteConfig
            {
                RouteId = "auth-admin",
                CorsPolicy = "WebOriginCorsPolicy",
                ClusterId = "auth-cluster",
                Match = new RouteMatch { Path = "/api/auth/admin/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/auth" } },
                ],
                AuthorizationPolicy = "admin-only",
            },
            new RouteConfig
            {
                RouteId = "bridge",
                ClusterId = "bridge-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/bridge/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/bridge" } },
                ],
                AuthorizationPolicy = "admin-or-peer",
            },
            new RouteConfig
            {
                RouteId = "stream",
                ClusterId = "stream-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/stream/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/stream" } },
                ],
                AuthorizationPolicy = "admin-or-peer",
            },
            new RouteConfig
            {
                RouteId = "resource1",
                ClusterId = "resource-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/resource/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/resource" } },
                ],
                AuthorizationPolicy = "admin-or-peer",
            },
            new RouteConfig
            {
                RouteId = "resource2",
                ClusterId = "resource-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/resource/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/resource" } },
                ],
                AuthorizationPolicy = "admin-or-peer",
            },
            new RouteConfig
            {
                RouteId = "resource3",
                ClusterId = "resource-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/resource/admin" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/resource" } },
                ],
                AuthorizationPolicy = "admin-only",
            },
            new RouteConfig
            {
                RouteId = "peer",
                ClusterId = "peer-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/peer/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/peer" } },
                ],
                AuthorizationPolicy = "admin-or-peer",
            },
        ];
    }
}
