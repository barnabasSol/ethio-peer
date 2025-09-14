using gateway.Service;
using gateway.YarpUtils;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Yarp.ReverseProxy.Configuration;

var builder = WebApplication.CreateBuilder(args);

var jwtSettings = builder.Configuration.GetSection("JwtSettings").Get<JwtSettings>();

builder.Services.AddSingleton<IJwksProvider>(new JwksProvider(jwtSettings!.JwksUri));

builder.Services.AddHostedService<JwksRefreshService>();

builder
    .Services.AddAuthentication(o =>
    {
        o.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
        o.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    })
    .AddJwtBearer();

builder.Services.AddCustomJwtAuthentication(jwtSettings);

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

app.MapControllers();

app.MapHealthChecks("/health");

app.UseAuthentication();
app.UseAuthorization();

app.UseCors("WebOriginCorsPolicy");

app.MapReverseProxy();

app.Run();
