using gateway.Service;
using gateway.YarpUtils;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;

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

builder
    .Services.AddOptions<JwtBearerOptions>(JwtBearerDefaults.AuthenticationScheme)
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
