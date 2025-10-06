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
                Match = new RouteMatch { Path = "/auth/{**catch-all}" },
                Transforms = [new Dictionary<string, string> { { "PathRemovePrefix", "/auth" } }],
                AuthorizationPolicy = "anonymous",
            },
            new RouteConfig
            {
                RouteId = "auth-reset-password",
                CorsPolicy = "WebOriginCorsPolicy",
                ClusterId = "auth-cluster",
                Match = new RouteMatch { Path = "/api/auth/password/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/auth" } },
                ],
                AuthorizationPolicy = "authenticated",
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
                AuthorizationPolicy = "authenticated",
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
                AuthorizationPolicy = "authenticated",
            },
            new RouteConfig
            {
                RouteId = "resource",
                ClusterId = "resource-cluster",
                CorsPolicy = "WebOriginCorsPolicy",
                Match = new RouteMatch { Path = "/api/resource/{**catch-all}" },
                Transforms =
                [
                    new Dictionary<string, string> { { "PathRemovePrefix", "/api/resource" } },
                ],
                AuthorizationPolicy = "authenticated",
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
                AuthorizationPolicy = "authenticated",
            },
        ];
    }
}
