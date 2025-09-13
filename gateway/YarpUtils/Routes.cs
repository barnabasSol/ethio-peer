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
                RouteId = "auth-reset-password",
                CorsPolicy = "WebOriginCorsPolicy",
                ClusterId = "auth-cluster",
                Match = new RouteMatch { Path = "/auth/health/{**catch-all}" },
                Transforms = [new Dictionary<string, string> { { "PathRemovePrefix", "/auth" } }],
                AuthorizationPolicy = "authenticated",
            },
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
        ];
    }
}
