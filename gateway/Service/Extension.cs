using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;

namespace gateway.Service;

public static class JwtAuthenticationExtensions
{
    public static IServiceCollection AddCustomJwtAuthentication(
        this IServiceCollection services,
        JwtSettings jwtSettings
    )
    {
        services
            .AddOptions<JwtBearerOptions>(JwtBearerDefaults.AuthenticationScheme)
            .Configure<IJwksProvider>(
                (options, jwksProvider) =>
                {
                    options.RequireHttpsMetadata = jwtSettings.RequireHttpsMetadata;

                    options.TokenValidationParameters = new TokenValidationParameters
                    {
                        ValidateIssuer = true,
                        ValidIssuer = jwtSettings.Issuer,

                        ValidateAudience = true,
                        ValidAudiences = jwtSettings.Audiences,

                        ValidateLifetime = true,
                        RequireExpirationTime = true,

                        ValidateIssuerSigningKey = true,

                        IssuerSigningKeyResolver = (
                            token,
                            securityToken,
                            kid,
                            validationParameters
                        ) =>
                        {
                            return jwksProvider.GetSigningKeys().Where(k => k.KeyId == kid) ?? [];
                        },
                    };

                    options.Events = new JwtBearerEvents
                    {
                        OnMessageReceived = ctx =>
                        {
                            if (
                                string.IsNullOrEmpty(ctx.Token)
                                && ctx.Request.Cookies.TryGetValue(
                                    "access_token",
                                    out var accessToken
                                )
                            )
                            {
                                ctx.Token = accessToken;
                            }

                            if (
                                ctx.Request.Cookies.TryGetValue(
                                    "refresh_token",
                                    out var refreshToken
                                )
                            )
                            {
                                ctx.HttpContext.Items["refresh_token"] = refreshToken;
                            }

                            return Task.CompletedTask;
                        },
                    };
                }
            );

        return services;
    }
}
