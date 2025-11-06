using System.Security.Claims;
using gateway.Service;
using gateway.YarpUtils;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Yarp.ReverseProxy.Configuration;
using System.Threading.RateLimiting;


var builder = WebApplication.CreateBuilder(args);

var jwtSettings = builder.Configuration.GetSection("JwtSettings").Get<JwtSettings>();

builder.Services.AddSingleton<IJwksProvider>(new JwksProvider(jwtSettings!.JwksUri));

builder.Services.AddHostedService<JwksRefreshService>();

builder.Services.AddRateLimiter(options =>
{
    options.GlobalLimiter = PartitionedRateLimiter.Create<HttpContext, string>(httpContext =>
    {
        var ip = httpContext.Connection.RemoteIpAddress?.ToString() ?? "unknown";

        return RateLimitPartition.GetFixedWindowLimiter(ip, _ => new FixedWindowRateLimiterOptions
        {
            PermitLimit = 100,
            Window = TimeSpan.FromMinutes(1),
            QueueProcessingOrder = QueueProcessingOrder.OldestFirst,
            QueueLimit = 0
        });
    });

    options.RejectionStatusCode = 429;
});

builder
.Services.AddAuthentication(o =>
{
    o.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
    o.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
})
.AddJwtBearer();

builder.Services.AddCustomJwtAuthentication(jwtSettings);

builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("admin-or-peer", policy =>
            policy.RequireRole("admin", "peer"));

    options.AddPolicy("admin-only", policy =>
        policy.RequireClaim(ClaimTypes.Role, "admin"));

    options.AddPolicy("peer-only", policy =>
        policy.RequireClaim(ClaimTypes.Role, "peer"));

    options.AddPolicy("Authenticated", policy =>
        policy.RequireAuthenticatedUser());
});

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
                .WithOrigins("http://localhost:5173", "https://ep-web.barney-host.site")
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
app.MapGet(
    "/proxy/config",
    (IProxyConfigProvider configProvider) =>
    {
        var config = configProvider.GetConfig();
        return Results.Json(config);
    }
);

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseCors("WebOriginCorsPolicy");

app.Use(
    async (context, next) =>
    {
        if (HttpMethods.IsOptions(context.Request.Method))
        {
            Console.WriteLine(
                $"[CORS Preflight] {context.Request.Method} {context.Request.Path} from {context.Request.Headers.Origin}"
            );
        }

        await next.Invoke();
    }
);

app.MapControllers();

app.MapHealthChecks("/health");

app.UseAuthentication();
app.UseAuthorization();

app.UseRateLimiter();

app.MapReverseProxy();

app.Run();
