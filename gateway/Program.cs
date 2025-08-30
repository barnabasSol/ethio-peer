using gateway.Service;
using gateway.YarpUtils;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using Yarp.ReverseProxy.Transforms;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<IJwksProvider>(
    new JwksProvider("http://localhost:2000/static/.well-known/jwks.json")
);

builder.Services.AddHostedService<JwksRefreshService>();

builder
    .Services.AddAuthentication(o =>
    {
        o.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
        o.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    })
    .AddJwtBearer();

builder
    .Services.AddOptions<JwtBearerOptions>(JwtBearerDefaults.AuthenticationScheme)
    .Configure<IJwksProvider>(
        (options, jwksProvider) =>
        {
            options.RequireHttpsMetadata = false;

            options.TokenValidationParameters = new TokenValidationParameters
            {
                ValidateIssuer = true,
                ValidIssuer = "ep.barney-host.site",

                ValidateAudience = true,
                ValidAudiences = ["https://ep-web.barney-host.site", "http://localhost:5173"],

                ValidateLifetime = true,
                RequireExpirationTime = true,

                ValidateIssuerSigningKey = true,

                IssuerSigningKeyResolver = (token, securityToken, kid, validationParameters) =>
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
                        && ctx.Request.Cookies.TryGetValue("access_token", out var accessToken)
                    )
                    {
                        ctx.Token = accessToken;
                    }

                    if (ctx.Request.Cookies.TryGetValue("refresh_token", out var refreshToken))
                    {
                        ctx.HttpContext.Items["refresh_token"] = refreshToken;
                    }

                    return Task.CompletedTask;
                },
            };
        }
    );

builder
    .Services.AddAuthorizationBuilder()
    .AddPolicy(
        "authenticated",
        policy =>
        {
            policy.RequireAuthenticatedUser();
        }
    );

builder.Services.AddCors(options =>
{
    options.AddPolicy(
        "WebOriginCorsPolicy",
        policy =>
        {
            policy
                .WithOrigins("http://localhost:5173")
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials();
        }
    );
});

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddHealthChecks();
builder.Services.AddSwaggerGen();

builder.Services.AddControllers();

builder
    .Services.AddReverseProxy()
    .LoadFromMemory(Routes.GetRoutes(), Clusters.GetClusters())
    .AddTransforms(builderContext =>
    {
        builderContext.RequestTransforms.Add(new ClaimsTransform());
    });

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.MapControllers();

app.MapHealthChecks("/health");

app.UseAuthentication();
app.UseAuthorization();

app.UseCors("WebOriginCorsPolicy");

app.MapReverseProxy();

app.Run();
