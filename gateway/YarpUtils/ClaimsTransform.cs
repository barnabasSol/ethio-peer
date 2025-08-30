using System.Security.Claims;
using Yarp.ReverseProxy.Transforms;

namespace gateway.YarpUtils;

public class ClaimsTransform : RequestTransform
{
    public override ValueTask ApplyAsync(RequestTransformContext context)
    {
        var user = context.HttpContext.User;
        if (user.Identity?.IsAuthenticated != true)
            return ValueTask.CompletedTask;

        // sub → NameIdentifier
        var sub = user.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (!string.IsNullOrEmpty(sub))
        {
            context.ProxyRequest.Headers.Remove("X-Claim-Sub");
            context.ProxyRequest.Headers.Add("X-Claim-Sub", sub);
        }

        var username = user.FindFirst("username")?.Value;
        if (!string.IsNullOrEmpty(username))
        {
            context.ProxyRequest.Headers.Remove("X-Claim-Username");
            context.ProxyRequest.Headers.Add("X-Claim-Username", username);
        }

        var jti = user.FindFirst("jti")?.Value;
        if (!string.IsNullOrEmpty(jti))
        {
            context.ProxyRequest.Headers.Remove("X-Claim-Jti");
            context.ProxyRequest.Headers.Add("X-Claim-Jti", jti);
        }
        var roles = user
            .Claims.Where(c => c.Type == ClaimTypes.Role)
            .Select(c => c.Value)
            .ToArray();

        if (roles.Length > 0)
        {
            context.ProxyRequest.Headers.Remove("X-Claim-Roles");
            context.ProxyRequest.Headers.Add("X-Claim-Roles", string.Join(",", roles));
        }

        return ValueTask.CompletedTask;
    }
}
